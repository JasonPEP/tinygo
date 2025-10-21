package http

import (
	stdhttp "net/http"
	"time"

	"tinygo/internal/logger"
	"tinygo/internal/shortener"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// NewMux creates a new mux router with all routes and middlewares
func NewMux(svc *shortener.Service) *mux.Router {
	handlers := NewHandlers(svc)

	// Create main router
	r := mux.NewRouter()

	// Management API routes (for admin/management)
	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/shorten", handlers.shorten).Methods("POST")
	admin.HandleFunc("/links/{code}", handlers.linkDetail).Methods("GET", "DELETE")
	admin.HandleFunc("/stats", handlers.stats).Methods("GET")

	// Public API routes (for programmatic access)
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/shorten", handlers.shorten).Methods("POST")
	api.HandleFunc("/links/{code}", handlers.linkDetail).Methods("GET", "DELETE")

	// Static files
	r.PathPrefix("/static/").Handler(stdhttp.StripPrefix("/static/", stdhttp.FileServer(stdhttp.Dir("web/static/"))))

	// Web UI
	r.HandleFunc("/", handlers.webUI).Methods("GET")

	// Health check endpoints
	r.HandleFunc("/healthz", handlers.health).Methods("GET")
	r.HandleFunc("/readyz", handlers.ready).Methods("GET")

	// Core feature: Short URL redirect (must be last to avoid conflicts)
	// This is the main purpose: ultra-short URLs like /abc123
	// Use a more specific matcher to avoid conflicts
	r.Path("/{code}").HandlerFunc(handlers.redirect).Methods("GET")

	// Apply middlewares
	r.Use(loggingMiddleware)
	r.Use(recoveryMiddleware)
	r.Use(corsMiddleware)

	return r
}

// Middleware functions
func loggingMiddleware(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		logger.Log.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start),
			"remote":   r.RemoteAddr,
		}).Info("request completed")
	})
}

func recoveryMiddleware(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log.WithFields(logrus.Fields{
					"panic":  rec,
					"path":   r.URL.Path,
					"method": r.Method,
				}).Error("panic recovered")
				w.WriteHeader(stdhttp.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(stdhttp.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
