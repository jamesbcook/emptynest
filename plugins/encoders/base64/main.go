package main

import (
	"encoding/base64"
)

func Name() string {
	return "base64"
}

func Encode(data []byte) ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(data)), nil
}

func Decode(data []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(data))
}
