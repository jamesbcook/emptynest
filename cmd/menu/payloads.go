package main

import (
	"fmt"
	"strconv"

	"github.com/empty-nest/server"
)

// TODO: print usage documentation
func (m *menu) payloads(args []string) {
	if len(args) < 1 {
		var payloads []emptynest.Payload
		m.DB.All(&payloads)
		for _, p := range payloads {
			fmt.Printf("%d %s\n", p.ID, p.Name)
		}
		return
	}

	switch args[0] {
	case "list":
		fmt.Printf("Available payloads:\n")
		for k := range m.PayloadMap {
			fmt.Printf("* %s\n", k)
		}
		return
	case "add":
		if len(args) < 3 {
			fmt.Println("syntax: payloads add [type] <name> <arguments>")
			return
		}
		plug, ok := m.PayloadMap[args[1]]
		if !ok {
			fmt.Printf("[!] payload plugin not found\n")
			return
		}
		name := args[2]
		data, err := plug.Generate(args[3:])
		if err != nil {
			fmt.Printf("[!] Error generating payload. Error %s\n", err.Error())
			return
		}
		payload := emptynest.Payload{
			Name: name,
			Kind: plug.Name(),
			Data: append([]byte{byte(plug.ID())}, data...),
		}
		m.DB.Save(&payload)
	case "del":
		if len(args) < 2 {
			fmt.Println("syntax: keys del <id>")
			return
		}
		keyID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("[!] There was an error parsing the id provided\n")
		}
		m.DB.Remove(&emptynest.Payload{ID: keyID})
	case "get":
		if len(args) < 2 {
			fmt.Println("syntax: payloads get <id>")
			return
		}
		keyID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("[!] There was an error parsing the id provided\n")
		}
		var payload emptynest.Payload
		if err := m.DB.One("ID", keyID, &payload); err != nil {
			fmt.Printf("[!] There was an error finding payload id %d. Error: %s\n", keyID, err.Error())
			return
		}
		plug := m.PayloadMap[payload.Kind]
		fmt.Printf("Name: %s Kind: %s Data: %s", payload.Kind, payload.Name, plug.String(payload.Data[1:]))
	}
}
