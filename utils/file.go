package utils

import (
	"fmt"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func GetMimetype(b []byte) string {
	mt := mimetype.Detect(b)
	return mt.String()
}

// returns contents, mimetype with charset, error
func ReadFile(path string) (string, string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file at '%s': %w", path, err)
	}
	return string(b), GetMimetype(b), nil
}

// returns contents, mimetype
func MustReadFile(path string) (string, string) {
	s, mt, err := ReadFile(path)
	if err != nil {
		panic(err)
	}
	return s, mt
}
