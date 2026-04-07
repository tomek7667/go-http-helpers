package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

func GenerateKeyPairRSA() (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	pubASN1 := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	return string(pubBytes), string(privBytes)
}

func getPubKey(pubKeyStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode public key PEM")
	}

	// Try PKCS1 first
	pubKey1, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return pubKey1, nil
	}
	pubKey2, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	rsaPubKey, ok := pubKey2.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return rsaPubKey, nil
}

// EncryptMessage encrypts a message using the recipient's public key without manual nonce handling.
// All keys are base64 encoded. The encrypted message is also base64 encoded.
func EncryptMessageRSA(message string, recipientPubKey string) (string, error) {
	pubKey, err := getPubKey(recipientPubKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptMessage decrypts a message using the recipient's private key
// All keys are base64 encoded, and the encrypted message is also base64 encoded.
func DecryptMessageRSA(encrypted string, recipientPrivKey string) (string, error) {
	privBlock, _ := pem.Decode([]byte(recipientPrivKey))
	if privBlock == nil {
		return "", fmt.Errorf("failed to decode private key PEM")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return "", err
	}
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, encryptedBytes)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func LongEncryptMessageRSA(message, recipientPubKey string) (string, error) {
	splitted := SplitMessage(message, 245)
	var encrypted []string
	for _, s := range splitted {
		enc, err := EncryptMessageRSA(s, recipientPubKey)
		if err != nil {
			return "", err
		}
		encrypted = append(
			encrypted,
			enc,
		)
	}
	return strings.Join(encrypted, ";"), nil
}

func LongDecryptMessageRSA(encrypted, recipientPrivKey string) (string, error) {
	splitted := strings.Split(encrypted, ";")
	decrypted := ""
	for _, s := range splitted {
		dec, err := DecryptMessageRSA(s, recipientPrivKey)
		if err != nil {
			return "", err
		}
		decrypted += dec
	}
	return decrypted, nil
}
