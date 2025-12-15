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
	// We allocate a buffer for the program counters (PCs)
	pc := make([]uintptr, 32)
	// Skip 3 frames: runtime.Callers, getCleanStack, ResErr
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	var builder strings.Builder

	for {
		frame, more := frames.Next()

		// --- FILTERING LOGIC ---

		// 1. Skip Go Runtime internals (runtime/...)
		if strings.Contains(frame.File, "runtime/") {
			if !more {
				break
			}
			continue
		}

		// 2. Stop when we hit the generic HTTP server guts
		// (This cuts off the massive bottom half of your original log)
		if strings.Contains(frame.Func.Name(), "net/http.HandlerFunc.ServeHTTP") ||
			strings.Contains(frame.Func.Name(), "net/http.serverHandler.ServeHTTP") {
			break
		}

		// 3. Optional: Make library paths shorter (remove version hash noise)
		// e.g. "go/pkg/mod/github.com/..." -> "github.com/..."
		cleanFile := frame.File
		if idx := strings.Index(frame.File, "go/pkg/mod/"); idx != -1 {
			cleanFile = frame.File[idx+11:] // Keep path after mod/
		} else {
			// Try to make project files relative
			if wd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(wd, frame.File); err == nil {
					cleanFile = rel
				}
			}
		}

		// Append to output
		builder.WriteString(fmt.Sprintf("\n\t%s\n\t\t%s:%d", frame.Func, cleanFile, frame.Line))

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
