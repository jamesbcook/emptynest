package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"plugin"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/empty-nest/emptynest"
)

func checkAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	db                  *storm.DB
	payloadMap          map[string]emptynest.PayloadPlugin
	transportMap        = map[string]*emptynest.Transport{}
	hostChanMap         = map[int]chan emptynest.ApprovalResponse{}
	approvalRequestChan = make(chan emptynest.ApprovalRequest)
	debugChan           = make(chan string)
	logChan             = make(chan string)
)

func main() {
	conf, err := decodeConfigFile("config.toml")
	checkAndPanic(err)

	db, err = storm.Open(conf.DBFile)
	checkAndPanic(err)
	err = db.Init(&emptynest.Host{})
	checkAndPanic(err)
	err = db.Init(&emptynest.Payload{})
	checkAndPanic(err)
	reader := bufio.NewReader(os.Stdin)
	payloadMap, err = emptynest.PayloadMap(conf.PayloadPluginDirectories)
	checkAndPanic(err)

	for _, trconf := range conf.Transports {
		p, err := plugin.Open(trconf.PluginLocation)
		if err != nil {
			fmt.Printf("[!] Could not open transport plugin at %s\n", trconf.PluginLocation)
			os.Exit(1)
		}
		createfunc, err := p.Lookup("Create")
		if err != nil {
			fmt.Println("[!] Transport plugin does not expose Create function")
			os.Exit(1)
		}
		encoderChain, err := emptynest.BuildEncoderChain(trconf.EncoderPluginLocations)
		checkAndPanic(err)
		cryptoChain, err := emptynest.BuildCryptoChain(trconf.CryptoPluginLocations)
		checkAndPanic(err)
		var keyChain [][]byte
		for _, k := range trconf.KeyChain {
			key, err := hex.DecodeString(k)
			checkAndPanic(err)
			keyChain = append(keyChain, key)
		}
		infoPlugin, err := emptynest.BuildHostInfoPlugin(trconf.HostInfoPluginLocation)
		checkAndPanic(err)
		create := createfunc.(func(emptynest.TransportCtx) (emptynest.Transport, error))
		transport, err := create(emptynest.TransportCtx{
			Addr:                trconf.Addr,
			ConfigFileLocation:  trconf.ConfigFileLocation,
			ApprovalRequestChan: approvalRequestChan,
			DebugChan:           debugChan,
			LogChan:             logChan,
			Debug:               conf.Debug,
			EncoderChain:        encoderChain,
			CryptoChain:         cryptoChain,
			HostInfo:            infoPlugin,
			KeyChain:            keyChain,
		})
		if err != nil {
			fmt.Println("[!] Error creating transport")
			os.Exit(1)
		}
		transportMap[transport.Name()] = &transport
		go func(tr emptynest.Transport, c emptynest.TransportConfig) {
			fmt.Printf("[+] Starting transport %s on %s\n", tr.Name(), c.Addr)
			if err := tr.Start(); err != nil {
				fmt.Printf("[!] Error starting trasnport %s\n", tr.Name())
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}
		}(transport, trconf)
	}

	// Listen for approval requests from the server.
	go func() {
		for {
			request := <-approvalRequestChan
			var (
				host   = request.Host
				exHost emptynest.Host
			)
			if err := db.Select(q.Eq("Info", host.Info)).First(&exHost); err == nil {
				host.ID = exHost.ID
				db.Update(&host)
			} else {
				db.Save(&host)
			}
			hostChanMap[host.ID] = request.Chan
			fmt.Printf("[!] APPROVAL REQUESTED:\nID: %d\nInfo:\n%s\n", host.ID, host.Info)
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
