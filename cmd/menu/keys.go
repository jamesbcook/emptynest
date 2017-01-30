package main

import (
	"encoding/hex"
	"fmt"

	"github.com/tomsteele/emptynest"
)

func (m *menu) keys(args []string) {
	if len(args) < 1 {
		var keys []emptynest.Key
		m.DB.All(&keys)
		for _, k := range keys {
			fmt.Printf("id: %d key: %s\n", k.ID, hex.EncodeToString(k.Key))
		}
		return
	}
	fmt.Println("not implemented.")
	return
}
