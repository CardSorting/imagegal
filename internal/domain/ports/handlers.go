package ports

import "net/http"

// Text2ImgHandler defines the interface for handling text-to-image requests
type Text2ImgHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// HealthHandler defines the interface for handling health check requests
type HealthHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// Handler represents a generic HTTP handler interface
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
