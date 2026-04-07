package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// deriveKey generates a 32-byte key from a string using SHA-512
func deriveKey(keyString string) []byte {
	hash := NfoldSha512(keyString, 5)
	return []byte(hash)[:32] // Use only the first 32 bytes for AES-256
}

// EncryptAES256 encrypts a string using AES-256-GCM with a SHA-512-derived key
func EncryptAES256(plaintext, key string) string {
	block, err := aes.NewCipher(deriveKey(key))
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, 12) // AES-GCM recommended nonce size
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
}

// DecryptAES256 decrypts a string using AES-256-GCM with a SHA-512-derived key
func DecryptAES256(ciphertext, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(deriveKey(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := data[:12]  // Extract nonce
	encMsg := data[12:] // Extract ciphertext

	plaintext, err := aesGCM.Open(nil, nonce, encMsg, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
