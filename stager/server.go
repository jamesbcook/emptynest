package stager

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	emptynest "github.com/tomsteele/emptynest"
	"github.com/tomsteele/emptynest/crypter"
)

// ApprovalRequest ...
type ApprovalRequest struct {
	Host emptynest.Host
	Chan chan emptynest.Payload
}

// Server is an HTTP transport to handle approval
// and denial of incoming stage requests.
type Server struct {
	ApproveRequest   chan ApprovalRequest
	DataLocation     string // Where in the in the request the DataParam will be.
	DataParam        string // Name of the parameter holding the Host information.
	CryptoByteOffset int    // Location of the encryption algorithm byte.
	BodyByteOffset   int    // Starting byte of the payload body.
	DB               *storm.DB
	Debug            bool
}

func (s *Server) debug(msg string) {
	if s.Debug {
		log.Println(msg)
	}
}

// Handle is the main HandlerFunc for handling incoming stage requests.
func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	payload, err := s.extractPayload(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if payload == "" {
		http.NotFound(w, r)
		return
	}
	decodeType, data, err := decodePayload(payload)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if len(data) < s.CryptoByteOffset+1 || len(data) < s.BodyByteOffset {
		http.NotFound(w, r)
		return
	}
	cryptoByte := data[s.CryptoByteOffset]
	s.debug(fmt.Sprintf("Cryptobyte: %v", cryptoByte))
	body := data[s.BodyByteOffset:]
	keyID := data[s.CryptoByteOffset+1]
	s.debug(fmt.Sprintf("Keybyte: %v", keyID))
	var key emptynest.Key
	if int(keyID) != 0 {
		if err := s.DB.One("ID", int(keyID), &key); err != nil {
			http.NotFound(w, r)
			return
		}
	}
	crypt, ok := crypter.Map[int(cryptoByte)]
	if !ok {
		http.NotFound(w, r)
		return
	}
	plain, err := crypt.Unseal(key.Key, body)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	parts := bytes.Split(plain, []byte{0xff})
	if len(parts) < 2 {
		s.debug("not enough info in payload")
		http.NotFound(w, r)
		return
	}
	var (
		host      emptynest.Host
		hostname  = string(parts[0])
		username  = string(parts[1])
		ipAddress = r.RemoteAddr
	)

	err = s.DB.Select(q.And(q.Eq("Hostname", hostname), q.Eq("Username", username))).First(&host)
	if err != nil {
		host = emptynest.Host{
			Hostname:  hostname,
			Username:  username,
			IPAddress: ipAddress,
		}

		if len(parts) > 2 {
			host.Misc = string(parts[2])
		}
		s.DB.Save(&host)
	}
	payloadChannel := make(chan emptynest.Payload)
	defer close(payloadChannel)

	s.ApproveRequest <- ApprovalRequest{Host: host, Chan: payloadChannel}
	response := <-payloadChannel
	if response.ID == 0 {
		http.NotFound(w, r)
		return
	}
	ciphertext, err := crypt.Seal(key.Key, response.Data)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	encodedBody, err := encodePayload(decodeType, ciphertext)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, encodedBody)
	return
}

func (s *Server) extractPayload(r *http.Request) (string, error) {
	var payload string
	var err error
	switch s.DataLocation {
	case "query":
		payload = r.URL.Query().Get(s.DataParam)
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

// Constant variables for encoding types.
const (
	BASE64 = 0
	HEX    = 1
)

// For now this will attempt base64 and then hex.
// Should validate by regex in the future
func decodePayload(payload string) (int, []byte, error) {
	if data, err := hex.DecodeString(payload); err == nil {
		return HEX, data, nil
	}
	data, err := base64.StdEncoding.DecodeString(payload)
	return BASE64, data, err
}

func encodePayload(encodeType int, data []byte) (string, error) {
	switch encodeType {
	case BASE64:
		return base64.StdEncoding.EncodeToString(data), nil
	case HEX:
		return hex.EncodeToString(data), nil
	default:
		return "", errors.New("not implemented")
	}
}
