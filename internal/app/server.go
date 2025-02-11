package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"image/internal/domain/ports"
	"image/internal/infrastructure/config"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	router *mux.Router
	logger ports.Logger
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, logger ports.Logger, handlers map[string]ports.Handler) *Server {
	router := mux.NewRouter()

	// Create server instance
	srv := &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
		router: router,
		logger: logger,
	}

	// Setup routes
	srv.setupRoutes(handlers)

	return srv
}

// setupRoutes configures the server routes
func (s *Server) setupRoutes(handlers map[string]ports.Handler) {
	// Add middleware to all routes - order matters!
	s.router.Use(s.corsMiddleware) // CORS headers must be first
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.recoveryMiddleware)

	// API routes
	api := s.router.PathPrefix("/api/v6").Subrouter()

	// Text to Image endpoint
	if h, ok := handlers["text2img"]; ok {
		api.Handle("/images/text2img", s.middleware(h)).Methods(http.MethodPost, http.MethodOptions)
	}

	// Health check endpoint
	if h, ok := handlers["health"]; ok {
		s.router.Handle("/health", h).Methods(http.MethodGet)
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting server", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")
	return s.server.Shutdown(ctx)
}

// middleware wraps a handler with common middleware
func (s *Server) middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

// corsMiddleware handles CORS headers for all routes
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs incoming requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		s.logger.Info("Request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)

		s.logger.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

// recoveryMiddleware recovers from panics
func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Error("Panic recovered", fmt.Errorf("%v", err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
