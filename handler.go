package emptynest

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

// ApprovalRequest ...
type ApprovalRequest struct {
	Host Host
	Chan chan ApprovalResponse
}

// ApprovalResponse ...
type ApprovalResponse struct {
	Payload Payload
	Plugin  PayloadPlugin
}

// Server is an HTTP transport to handle approval
// and denial of incoming stage requests.
type Server struct {
	ApproveRequest chan ApprovalRequest
	DB             *storm.DB
	EncoderChain   []EncoderPlugin
	CryptoChain    []CryptoPlugin
	KeyChain       [][]byte
	HostInfo       HostInfoPlugin
	GetLocation    string
	GetParam       string
	PostLocation   string
	PostParam      string
	Debug          bool
}

func (s *Server) debug(msg string) {
	if s.Debug {
		log.Println(msg)
	}
}

// Handle is the main HandlerFunc for handling incoming stage requests.
func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	var (
		payload string
		err     error
	)
	switch r.Method {
	case "GET":
		payload, err = s.extractGET(r)
	case "POST":
		payload, err = s.extractPOST(r)
	default:
		err = errors.New("not implemented")
	}
	if payload == "" || err != nil {
		http.NotFound(w, r)
		return
	}

	data := []byte(payload)
	for _, encoder := range s.EncoderChain {
		s.debug("Decoding with " + encoder.Name())
		data, err = encoder.Decode(data)
		if err != nil {
			http.NotFound(w, r)
			s.debug(err.Error())
			return
		}
	}
	var key = s.KeyChain[0]
	for i, crypter := range s.CryptoChain {
		// If the length of keys is the length of plugins
		// then iterate over the keys. Else use the same key.
		if len(s.KeyChain) == len(s.CryptoChain) {
			key = s.KeyChain[i]
		}
		s.debug("Decrypting with " + crypter.Name())
		data, err = crypter.Open(key, data)
		if err != nil {
			s.debug(err.Error())
			http.NotFound(w, r)
			return
		}
	}
	parts := bytes.Split(data, s.HostInfo.SplitPattern())
	argLength := s.HostInfo.ArgLength()
	if len(parts) < argLength {
		s.debug("incorrect info length")
		http.NotFound(w, r)
		return
	}
	info := parts[:argLength]
	host := Host{
		IPAddress: r.RemoteAddr,
		Info:      s.HostInfo.String(info),
	}
	err = s.DB.Select(q.Eq("Info", info)).First(&host)
	if len(parts) > argLength {
		host.Data = string(parts[argLength+1])
	}
	s.DB.Save(&host)

	approvalResponseChan := make(chan ApprovalResponse)
	defer close(approvalResponseChan)
	s.ApproveRequest <- ApprovalRequest{
		Host: host,
		Chan: approvalResponseChan,
	}
	approvalResponse := <-approvalResponseChan
	if approvalResponse.Payload.ID == 0 {
		http.NotFound(w, r)
		return
	}
	response, err := approvalResponse.Plugin.Generate(approvalResponse.Payload.Data)
	if err != nil {
		s.debug(err.Error())
		http.NotFound(w, r)
		return
	}
	key = s.KeyChain[0]
	for i := (len(s.CryptoChain) - 1); i >= 0; i-- {
		if len(s.KeyChain) == len(s.CryptoChain) {
			key = s.KeyChain[i]
		}
		crypter := s.CryptoChain[i]
		s.debug("Encrypting with " + crypter.Name())
		response, err = crypter.Seal(key, response)
		if err != nil {
			s.debug(err.Error())
			http.NotFound(w, r)
			return
		}
	}
	for i := (len(s.EncoderChain) - 1); i >= 0; i-- {
		encoder := s.EncoderChain[i]
		s.debug("Encoding with " + encoder.Name())
		response, err = encoder.Encode(response)
		if err != nil {
			s.debug(err.Error())
			http.NotFound(w, r)
			return
		}
	}
	fmt.Fprint(w, string(response))
	return
}

func (s *Server) extractGET(r *http.Request) (string, error) {
	var payload string
	var err error
	switch s.GetLocation {
	case "query":
		payload = r.URL.Query().Get(s.GetParam)
	case "cookie":
		err = errors.New("not implemented")
	case "header":
		err = errors.New("not implemented")
	case "body":
		err = errors.New("not implemented")
	default:
		err = errors.New("invalid DataLocation for server")
	}
	return payload, err
}

func (s *Server) extractPOST(r *http.Request) (string, error) {
	var payload string
	var err error
	switch s.PostLocation {
	case "body":
		payload = r.FormValue(s.PostParam)
	default:
		err = errors.New("invalid DataLocation for server")
	}
	return payload, err
}
