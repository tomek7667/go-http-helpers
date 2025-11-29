package h

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetResponse[T any](response *http.Response) (*T, error) {
	var rb T
	err := json.NewDecoder(response.Body).Decode(&rb)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return &rb, nil
}

func MustGetResponse[T any](response *http.Response) *T {
	rb, err := GetResponse[T](response)
	if err != nil {
		panic(err)
	}
	return rb
}
