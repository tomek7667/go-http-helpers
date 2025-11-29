package utils

import "encoding/json"

func MustMarshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func MustMarshalB(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func MustUnmarshalB[T any](b []byte) T {
	var t T
	err := json.Unmarshal(b, &t)
	if err != nil {
		panic(err)
	}
	return t
}
