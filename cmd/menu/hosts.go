package main

import (
	"fmt"
	"strconv"

	"github.com/empty-nest/server"
)

func (m *menu) hosts(args []string) {
	if len(args) < 1 {
		var hosts []emptynest.Host
		m.DB.All(&hosts)
		for _, h := range hosts {
			fmt.Printf("%d %s %s %s %s\n", h.ID, h.IPAddress, h.Hostname, h.Username, h.Status)
			if h.Misc != "" {
				fmt.Printf("\t%s\n", h.Misc)
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
		var h emptynest.Host
		keyID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("[!] There was an error parsing the id provided\n")
			return
		}
		if err := m.DB.One("ID", keyID, &h); err != nil {
			fmt.Printf("[!] There was an error finding key id %d. Error: %s\n", keyID, err.Error())
			return
		}
		fmt.Printf("%d %s %s %s\n", h.ID, h.IPAddress, h.Hostname, h.Username)
		if h.Misc != "" {
			fmt.Printf("\t%s\n", h.Misc)
		}

	case "del":
		if len(args) < 2 {
			fmt.Println("syntax: keys del <id>")
			return
		}
		keyID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("[!] There was an error parsing the id provided\n")
		}
		m.DB.Remove(&emptynest.Host{ID: keyID})
	case "approve":
		if len(args) < 3 {
			fmt.Println("syntax: hosts approve <id> <payload_name>")
			return
		}
		var h emptynest.Host
		keyID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("[!] There was an error parsing the id provided\n")
			return
		}
		if err := m.DB.One("ID", keyID, &h); err != nil {
			fmt.Printf("[!] There was an error finding host id %d. Error: %s\n", keyID, err.Error())
			return
		}
		var payload emptynest.Payload
		if err := m.DB.One("Name", args[2], &payload); err != nil {
			fmt.Printf("[!] There was an error finding payload name %s. Error: %s\n", args[2], err.Error())
			return
		}
		h.Status = "approved"
		c, ok := m.HostChanMap[keyID]
		if !ok {
			fmt.Printf("[!] There was an error finding host id %d.\n", keyID)
			return
		}
		c <- payload
		m.DB.Save(&h)
		delete(m.HostChanMap, keyID)
		fmt.Printf("[+] Host %d approved\n", h.ID)
	case "deny":
		if len(args) < 2 {
			fmt.Println("syntax: hosts deny <id>")
			return
		}
		var h emptynest.Host
		keyID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("[!] There was an error parsing the id provided\n")
			return
		}
		if err := m.DB.One("ID", keyID, &h); err != nil {
			fmt.Printf("[!] There was an error finding host id %d. Error: %s\n", keyID, err.Error())
			return
		}
		h.Status = "denied"
		c, ok := m.HostChanMap[keyID]
		if !ok {
			fmt.Printf("[!] There was an error finding host id %d.\n", keyID)
			return
		}
		var payload emptynest.Payload
		c <- payload
		m.DB.Save(&h)
		delete(m.HostChanMap, keyID)
		fmt.Printf("[+] Host %d denied\n", h.ID)
	}
}
