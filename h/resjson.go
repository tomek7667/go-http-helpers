package h

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

func ResErr(w http.ResponseWriter, err error) {
	_, file, line, _ := runtime.Caller(1)
	slog.Error(
		"someone got an error from api",
		"err", err,
		"source", fmt.Sprintf("%s:%d", file, line),
	)
	slog.Debug(
		"someone got an error from api- debug stack tree",
		"err", err,
		"trace", string(debug.Stack()),
	)
	resjson(w, nil, err.Error(), http.StatusInternalServerError)
}

func ResNotFound(w http.ResponseWriter, resourceType string) {
	resjson(w, nil, fmt.Sprintf("%s was not found", resourceType), http.StatusNotFound)
}

func ResBadRequest(w http.ResponseWriter, err error) {
	resjson(w, nil, err.Error(), http.StatusBadRequest)
}

func ResSuccess(w http.ResponseWriter, data any) {
	resjson(w, data, "Success", http.StatusOK)
}

func ResUnauthorized(w http.ResponseWriter) {
	resjson(w, nil, "unauthorized", http.StatusUnauthorized)
}

func resjson(w http.ResponseWriter, data any, message string, code int) {
	marshallable := map[string]any{
		"success":    code == http.StatusOK,
		"message":    message,
		"data":       data,
		"statusText": http.StatusText(code),
		"code":       code,
		"ts":         time.Now(),
	}
	b, _ := json.Marshal(marshallable)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}
