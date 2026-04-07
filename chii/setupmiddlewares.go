package chii

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

func isDebugHeaderSet(r *http.Request) bool {
	return r.Header.Get("Debug") == "reveal-body-logs"
}

func SetupMiddlewares(router chi.Router, allowedOrigins []string) {
	router.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{
		// Level defines the verbosity of the request logs:
		// slog.LevelDebug - log all responses (incl. OPTIONS)
		// slog.LevelInfo  - log responses (excl. OPTIONS)
		// slog.LevelWarn  - log 4xx and 5xx responses only (except for 429)
		// slog.LevelError - log 5xx responses only
		Level:  slog.LevelDebug,
		Schema: httplog.SchemaOTEL,

		// RecoverPanics recovers from panics occurring in the underlying HTTP handlers
		// and middlewares. It returns HTTP 500 unless response status was already set.
		//
		// NOTE: Panics are logged as errors automatically, regardless of this setting.
		RecoverPanics: true,

		// Optionally, filter out some request logs.
		Skip: func(req *http.Request, respStatus int) bool {
			return respStatus == 404 || respStatus == 405
		},

		// Optionally, log selected request/response headers explicitly.
		LogRequestHeaders:  []string{"Origin", "CF-Connecting-IP"},
		LogResponseHeaders: []string{},

		// Optionally, enable logging of request/response body based on custom conditions.
		// Useful for debugging payload issues in development.
		LogRequestBody:  isDebugHeaderSet,
		LogResponseBody: isDebugHeaderSet,
	}))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Use(middleware.Timeout(30 * time.Minute))
}
