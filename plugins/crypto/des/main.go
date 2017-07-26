package main

import (
	"crypto/cipher"
	"crypto/des"
)

func Name() string {
	return "3DES"
}

func Open(key, ciphertext []byte) ([]byte, error) {
	var tripleDESKey []byte
	tripleDESKey = append(tripleDESKey, key[:16]...)
	tripleDESKey = append(tripleDESKey, key[:8]...)
	block, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		return ciphertext, err
	}
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[des.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

func Seal(key, plaintext []byte) ([]byte, error) {
	var tripleDESKey []byte
	tripleDESKey = append(tripleDESKey, key[:16]...)
	tripleDESKey = append(tripleDESKey, key[:8]...)
	block, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		return plaintext, err
	}
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCTR(block, key[des.BlockSize:])
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}
