package h

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func GetStringResponse(resp *http.Response) (string, error) {
	if resp == nil || resp.Body == nil {
		return "", nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("io.ReadAll: %w", err)
	}

	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(b))

	return string(b), nil
}

func MustGetStringResponse(resp *http.Response) string {
	s, err := GetStringResponse(resp)
	if err != nil {
		panic(err)
	}
	return s
}
