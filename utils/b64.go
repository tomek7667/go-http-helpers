package utils

import (
	"encoding/base64"
	"fmt"
)

func B64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func B64Decode(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("failed to b64 decode '%s': %w", s, err)
	}
	return string(b), nil
}

func MustB64Decode(s string) string {
	s, err := B64Decode(s)
	if err != nil {
		panic(err)
	}
	return s
}
