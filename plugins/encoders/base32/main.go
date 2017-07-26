package main

import (
	"encoding/base32"
)

func Name() string {
	return "base32"
}

func Encode(data []byte) ([]byte, error) {
	return []byte(base32.StdEncoding.EncodeToString(data)), nil
}

func Decode(data []byte) ([]byte, error) {
	return base32.StdEncoding.DecodeString(string(data))
}
