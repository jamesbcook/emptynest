package main

import (
	"fmt"
	"strconv"

	"github.com/empty-nest/emptynest"
)

func getHost(id string) {
	var h emptynest.Host
	keyID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("[!] There was an error parsing the id provided\n")
		return
	}
	if err := db.One("ID", keyID, &h); err != nil {
		fmt.Printf("[!] There was an error finding key id %d. Error: %s\n", keyID, err.Error())
		return
	}
	fmt.Printf("%d %s\n", h.ID, h.Info)
	if h.Data != "" {
		fmt.Printf("\t%s\n", h.Data)
	}
}

func delHost(id string) {
	keyID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("[!] There was an error parsing the id provided\n")
	}
	db.Remove(&emptynest.Host{ID: keyID})
}

func approveHost(args []string) {
	var h emptynest.Host
	keyID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("[!] There was an error parsing the id provided\n")
		return
	}
	if err := db.One("ID", keyID, &h); err != nil {
		fmt.Printf("[!] There was an error finding host id %d. Error: %s\n", keyID, err.Error())
		return
	}
	var payload emptynest.Payload
	if err := db.One("Name", args[1], &payload); err != nil {
		fmt.Printf("[!] There was an error finding payload name %s. Error: %s\n", args[1], err.Error())
		return
	}
	h.Status = "approved"
	c, ok := hostChanMap[keyID]
	if !ok {
		fmt.Printf("[!] There was an error finding host id %d.\n", keyID)
		return
	}
	plug := payloadMap[payload.Kind]
	c <- emptynest.ApprovalResponse{
		Payload: payload,
		Plugin:  plug,
	}
	db.Save(&h)
	delete(hostChanMap, keyID)
	fmt.Printf("[+] Host %d approved\n", h.ID)
}

func denyHost(id string) {
	var h emptynest.Host
	keyID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("[!] There was an error parsing the id provided\n")
		return
	}
	if err := db.One("ID", keyID, &h); err != nil {
		fmt.Printf("[!] There was an error finding host id %d. Error: %s\n", keyID, err.Error())
		return
	}
	h.Status = "denied"
	c, ok := hostChanMap[keyID]
	if !ok {
		fmt.Printf("[!] There was an error finding host id %d.\n", keyID)
		return
	}
	var payload emptynest.Payload
	c <- emptynest.ApprovalResponse{
		Payload: payload,
	}
	db.Save(&h)
	delete(hostChanMap, keyID)
	fmt.Printf("[+] Host %d denied\n", h.ID)
}

func hostsHelp() {
	fmt.Println(`
get	Get list of hosts that have connected
del     Delete host that have connected
approve Approve connected host
deny    Deny connected host
`)
}

func hosts(args []string) {
	if len(args) < 1 {
		var hosts []emptynest.Host
		db.All(&hosts)
		for _, h := range hosts {
			fmt.Printf("%d %s\n", h.ID, h.Info)
			if h.Data != "" {
				fmt.Printf("\t%s\n", h.Data)
			}
		}
		return
	}
	switch args[0] {
	case "get":
		if len(args) < 2 {
			fmt.Println("syntax: hosts get <id>")
			return
		}
		getHost(args[1])
	case "del":
		if len(args) < 2 {
			fmt.Println("syntax: hosts del <id>")
			return
		}
		delHost(args[1])
	case "approve":
		if len(args) < 3 {
			fmt.Println("syntax: hosts approve <id> <payload_name>")
			return
		}
		approveHost(args[1:])
	case "deny":
		if len(args) < 2 {
			fmt.Println("syntax: hosts deny <id>")
			return
		}
		denyHost(args[1])
	case "help":
		hostsHelp()
		return
	}
}
