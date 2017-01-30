package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/asdine/storm"
	"github.com/tomsteele/emptynest"
	"github.com/tomsteele/emptynest/stager"
)

func checkAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// HostChanMap ...
var HostChanMap = map[int]chan emptynest.Payload{}

func main() {
	db, err := storm.Open("data.db")
	checkAndPanic(err)
	err = db.Init(&emptynest.Key{})
	checkAndPanic(err)
	err = db.Init(&emptynest.Host{})
	checkAndPanic(err)
	err = db.Init(&emptynest.Payload{})
	checkAndPanic(err)

	reader := bufio.NewReader(os.Stdin)

	// TODO: parse server variables
	// TODO: catch ctrl-c db.Close()
	fmt.Println("[-] Starting server")
	serve := stager.Server{}
	serve.BodyByteOffset = 2
	serve.CryptoByteOffset = 0
	serve.DataLocation = "query"
	serve.DataParam = "JSESSIONID"
	serve.ApproveRequest = make(chan stager.ApprovalRequest)
	serve.DB = db

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
			HostChanMap[host.ID] = request.Chan
			fmt.Printf("[!] APPROVAL REQUESTED:\nID: %d\nHostname: %s\nUsername: %s\nIP: %s\n", host.ID, host.Hostname, host.Username, host.IPAddress)
			fmt.Print("$> ")
		}
	}()

	pmap, err := emptynest.PayloadMap("./plugins")
	checkAndPanic(err)
	m := menu{DB: db, HostChanMap: HostChanMap, PayloadMap: pmap}

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
		case "keys":
			m.keys(args)
		case "hosts":
			m.hosts(args)
		case "cryptors":
			fmt.Printf("* none\n")
		case "payloads":
			m.payloads(args)
		case "help":
			fmt.Printf("Valid comamnds:\nkeys\nhosts\ncryptors\npayloads\n")
		default:
			fmt.Printf("[!] Unknown command %s. Try 'help'\n", text)
		}
	}
}
