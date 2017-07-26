package main

import (
	"encoding/hex"
)

func Name() string {
	return "hex"
}

func Encode(data []byte) ([]byte, error) {
	return []byte(hex.EncodeToString(data)), nil
}

func Decode(data []byte) ([]byte, error) {
	return hex.DecodeString(string(data))
}
