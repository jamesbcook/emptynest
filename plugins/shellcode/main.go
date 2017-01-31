package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Name is a descriptive name of the
func Name() string {
	return "shellcode"
}

// ID returns a unique integer across plugins
func ID() int {
	return 1
}

// Help returns documentation for the plugin.
func Help() string {
	return "shellcode <file|shellcode>"
}

// String returns as string friendly version of the payload.
func String(data []byte) string {
	return fmt.Sprintf("%x", data)
}

// Generate ...
func Generate(data []byte) ([]byte, error) {
	return data, nil
}

// Process prepares a payload based on the arguments provided
func Process(args []string) ([]byte, error) {
	var shellcode []byte
	if len(args) < 1 {
		return shellcode, errors.New("missing required arguments")
	}
	if len(args) > 1 {
		return shellcode, errors.New("invalid arguments provided")
	}
	if _, err := os.Stat(args[0]); err == nil {
		shellcode, err = ioutil.ReadFile(args[0])
		if err != nil {
			return shellcode, fmt.Errorf("There was an error reading binary file. Error:%s", err.Error())
		}
		return shellcode, nil
	}
	shellcode, err := hex.DecodeString(strings.Replace(args[0], "\\x", "", -1))
	if err != nil {
		return shellcode, fmt.Errorf("There was an error decoding shellcode. Error:%s", err.Error())

	}
	return shellcode, nil
}
