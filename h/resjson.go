package h

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func ResErr(w http.ResponseWriter, err error) {
	// 1. Get the source location of the error (caller of ResErr)
	_, file, line, _ := runtime.Caller(1)

	// Make path relative to project root (remove C:/Users/...)
	if wd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(wd, file); err == nil {
			file = rel
		}
	}

	slog.Error(
		"someone got an error from api",
		"err", err,
		"source", fmt.Sprintf("%s:%d", file, line),
	)

	// 2. Generate a clean stack trace
	slog.Debug(
		"api error trace",
		"err", err,
		"trace", getCleanStack(),
	)

	resjson(w, nil, err.Error(), http.StatusInternalServerError)
}

// getCleanStack returns a formatted string of the stack trace,
// filtering out noisy Go runtime and server internals.
func getCleanStack() string {
	pc := make([]uintptr, 32)
	// Skip 2 frames: runtime.Callers + getCleanStack
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var builder strings.Builder

	for {
		frame, more := frames.Next()

		// 1. Filter out Go runtime internals
		if strings.Contains(frame.File, "runtime/") {
			if !more {
				break
			}
			continue
		}

		// 2. Stop when we hit the generic HTTP server guts
		// This prevents the massive log tail
		if strings.Contains(frame.Function, "net/http.HandlerFunc.ServeHTTP") ||
			strings.Contains(frame.Function, "net/http.serverHandler.ServeHTTP") {
			break
		}

		// 3. Cleanup File Path (Relative or shorten go mod paths)
		cleanFile := frame.File
		// Remove messy go/pkg/mod/ prefix
		if idx := strings.Index(cleanFile, "go/pkg/mod/"); idx != -1 {
			cleanFile = cleanFile[idx+11:]
		} else {
			// Make project files relative to current working directory
			if wd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(wd, cleanFile); err == nil {
					cleanFile = rel
				}
			}
		}

		// 4. Format the output
		// FIX: Use frame.Function (string) instead of frame.Func (pointer)
		builder.WriteString(fmt.Sprintf("\n%s\n\t%s:%d", frame.Function, cleanFile, frame.Line))

		if !more {
			break
		}
	}

	return builder.String()
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
