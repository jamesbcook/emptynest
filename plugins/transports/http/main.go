package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/empty-nest/emptynest"
)

type config struct {
	emptynest.TransportConfig
	GetLocation  string `toml:"get_location"`
	GetParam     string `toml:"get_param"`
	PostLocation string `toml:"post_location"`
	PostParam    string `toml:"post_param"`
}

// Server implements emptynest.Server.
type Server struct {
	Ctx          emptynest.TransportCtx
	GetLocation  string
	GetParam     string
	PostLocation string
	PostParam    string
}

func (s *Server) debug(msg string) {
	if s.Ctx.Debug {
		go func(msg string) {
			s.Ctx.DebugChan <- msg
		}(msg)
	}
}

func (s *Server) log(msg string) {
	go func(msg string) {
		s.Ctx.LogChan <- msg
	}(msg)
}

// Handle is the HandlerFunc for handling incoming requests.
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
	for _, encoder := range s.Ctx.EncoderChain {
		s.debug("Decoding with " + encoder.Name())
		data, err = encoder.Decode(data)
		if err != nil {
			http.NotFound(w, r)
			s.debug(err.Error())
			return
		}
	}
	var key = s.Ctx.KeyChain[0]
	for i, crypter := range s.Ctx.CryptoChain {
		// If the length of keys is the length of plugins
		// then iterate over the keys. Else use the same key.
		if len(s.Ctx.KeyChain) == len(s.Ctx.CryptoChain) {
			key = s.Ctx.KeyChain[i]
		}
		s.debug("Decrypting with " + crypter.Name())
		data, err = crypter.Open(key, data)
		if err != nil {
			s.debug(err.Error())
			http.NotFound(w, r)
			return
		}
	}
	parts := bytes.Split(data, s.Ctx.HostInfo.SplitPattern())
	argLength := s.Ctx.HostInfo.ArgLength()
	if len(parts) < argLength {
		s.debug("incorrect info length")
		http.NotFound(w, r)
		return
	}
	info := s.Ctx.HostInfo.String(parts[:argLength])

	host := emptynest.Host{
		IPAddress: r.RemoteAddr,
		Info:      info,
	}
	if len(parts) > argLength {
		host.Data = string(parts[argLength])
	}
	approvalResponseChan := make(chan emptynest.ApprovalResponse)
	defer close(approvalResponseChan)
	s.Ctx.ApprovalRequestChan <- emptynest.ApprovalRequest{
		Host: host,
		Chan: approvalResponseChan,
	}
	approvalResponse := <-approvalResponseChan
	if approvalResponse.Payload.ID == 0 {
		http.NotFound(w, r)
		return
	}
	generatedPayload, err := approvalResponse.Plugin.Generate(approvalResponse.Payload.Data)
	if err != nil {
		s.debug(err.Error())
		http.NotFound(w, r)
		return
	}
	response := append([]byte{byte(approvalResponse.Plugin.ID())}, generatedPayload...)
	key = s.Ctx.KeyChain[0]
	for i := (len(s.Ctx.CryptoChain) - 1); i >= 0; i-- {
		if len(s.Ctx.KeyChain) == len(s.Ctx.CryptoChain) {
			key = s.Ctx.KeyChain[i]
		}
		crypter := s.Ctx.CryptoChain[i]
		s.debug("Encrypting with " + crypter.Name())
		response, err = crypter.Seal(key, response)
		if err != nil {
			s.debug(err.Error())
			http.NotFound(w, r)
			return
		}
	}
	for i := (len(s.Ctx.EncoderChain) - 1); i >= 0; i-- {
		encoder := s.Ctx.EncoderChain[i]
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

// Name returns "HTTP"
func (s *Server) Name() string {
	return "HTTP"
}

// Start starts the server.
func (s *Server) Start() error {
	http.HandleFunc("/", s.Handle)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "PONG\n")
	})
	return http.ListenAndServe(s.Ctx.Addr, nil)
}

// Stop stops the server.
func (s *Server) Stop() error {
	return nil
}

// Create creates a Server.
func Create(ctx emptynest.TransportCtx) (emptynest.Transport, error) {
	var (
		conf   config
		server Server
	)
	_, err := toml.DecodeFile(ctx.ConfigFileLocation, &conf)
	if err != nil {
		return &server, err
	}
	return &Server{
		Ctx:          ctx,
		GetLocation:  conf.GetLocation,
		GetParam:     conf.GetParam,
		PostLocation: conf.PostLocation,
		PostParam:    conf.PostParam,
	}, nil
}
