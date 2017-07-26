package main

import (
	"crypto/rc4"
)

func Name() string {
	return "RC4"
}

func Open(key, ciphertext []byte) ([]byte, error) {
	out := make([]byte, len(ciphertext))
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return out, err
	}
	cipher.XORKeyStream(out, ciphertext)
	return out, nil
}

func Seal(key, plaintext []byte) ([]byte, error) {
	out := make([]byte, len(plaintext))
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return out, err
	}
	cipher.XORKeyStream(out, plaintext)
	return out, nil

}
