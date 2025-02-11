package ports

import (
	"image/internal/domain/models"
)

// ModelRegistry defines the interface for managing AI models
type ModelRegistry interface {
	// Register adds a new model to the registry
	Register(model models.AIModel) error
	// Get retrieves a model by its ID
	Get(id string) (models.AIModel, error)
	// List returns all registered models
	List() []models.AIModel
	// Validate checks if a model ID is valid
	Validate(id string) error
}
