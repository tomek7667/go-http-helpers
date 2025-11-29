package h

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetDto[T any](r *http.Request) (*T, error) {
	var dto T
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling the body failed: %w", err)
	}
	return &dto, nil
}
