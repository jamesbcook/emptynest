package main

import "fmt"

func ArgLength() int {
	return 2
}

func SplitPattern() []byte {
	return []byte{0xff}
}

func String(in [][]byte) string {
	var (
		username = string(in[0])
		hostname = string(in[1])
	)
	return fmt.Sprintf("Username: %s\nHostname: %s\n", username, hostname)
}
