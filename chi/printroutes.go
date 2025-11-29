package chi

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
)

func PrintRoutes(router chi.Router) {
	// Helper function to walk through all registered routes
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		slog.Info("registered route", "method", method, "path", route)
		return nil
	}

	// Walk through all registered routes
	if err := chi.Walk(router, walkFunc); err != nil {
		slog.Error("error walking routes", "error", err)
	}
}
