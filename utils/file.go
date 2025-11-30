package utils

import (
	"fmt"
	"os"
)

func ReadFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file at '%s': %w", path, err)
	}
	return string(b), nil
}

func MustReadFile(path string) string {
	s, err := ReadFile(path)
	if err != nil {
		panic(err)
	}
	return s
}
