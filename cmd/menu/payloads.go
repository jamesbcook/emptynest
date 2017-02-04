package main

import (
	"fmt"
	"strconv"

	"github.com/empty-nest/emptynest"
)

func listPayloads() {
	fmt.Printf("Available payload types:\n")
	for k := range payloadMap {
		fmt.Printf("* %s\n", k)
	}
	return
}

func addPayload(args []string) {
	plug, ok := payloadMap[args[0]]
	if !ok {
		fmt.Printf("[!] payload plugin not found\n")
		return
	}
	name := args[1]
	data, err := plug.Process(args[2:])
	if err != nil {
		fmt.Printf("[!] Error generating payload. Error %s\n", err.Error())
		return
	}
	payload := emptynest.Payload{
		Name: name,
		Kind: plug.Name(),
		Data: append([]byte{byte(plug.ID())}, data...),
	}
	db.Save(&payload)
}

func delPayload(id string) {
	keyID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("[!] There was an error parsing the id provided\n")
	}
	db.Remove(&emptynest.Payload{ID: keyID})
}

func getPayload(id string) {
	keyID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("[!] There was an error parsing the id provided\n")
	}
	var payload emptynest.Payload
	if err := db.One("ID", keyID, &payload); err != nil {
		fmt.Printf("[!] There was an error finding payload id %d. Error: %s\n", keyID, err.Error())
		return
	}
	plug := payloadMap[payload.Kind]
	fmt.Printf("Name: %s Kind: %s Data: %s", payload.Kind, payload.Name, plug.String(payload.Data[1:]))

}

func payloads(args []string) {
	if len(args) < 1 {
		var payloads []emptynest.Payload
		db.All(&payloads)
		for _, p := range payloads {
			fmt.Printf("%d %s\n", p.ID, p.Name)
		}
		return
	}
	switch args[0] {
	case "list":
		listPayloads()
	case "add":
		if len(args) < 3 {
			fmt.Println("syntax: payloads add [type] <name> <arguments>")
			return
		}
		addPayload(args[1:])
	case "del":
		if len(args) < 2 {
			fmt.Println("syntax: payloads del <id>")
			return
		}
		delPayload(args[1])
	case "get":
		if len(args) < 2 {
			fmt.Println("syntax: payloads get <id>")
			return
		}
		getPayload(args[1])
	}
}
