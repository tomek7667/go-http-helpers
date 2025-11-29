package h

import (
	"fmt"
	"io"
	"net/http"
)

func GetStringResponse(response *http.Response) (string, error) {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to io.ReadAll: %w", err)
	}
	return string(b), nil
}

func MustGetStringResponse(response *http.Response) string {
	s, err := GetStringResponse(response)
	if err != nil {
		panic(err)
	}
	return s
}
