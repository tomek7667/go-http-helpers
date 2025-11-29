package h

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetMapResponse(response *http.Response) (map[string]any, error) {
	var rb map[string]any
	err := json.NewDecoder(response.Body).Decode(&rb)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return rb, nil
}

func MustGetMapResponse(response *http.Response) map[string]any {
	rb, err := GetMapResponse(response)
	if err != nil {
		panic(err)
	}
	return rb
}
