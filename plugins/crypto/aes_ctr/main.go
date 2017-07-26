package main

import (
	"crypto/aes"
	"crypto/cipher"
)

func Name() string {
	return "AES"
}

func Open(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return ciphertext, err
	}
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

func Seal(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext, err
	}
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}
