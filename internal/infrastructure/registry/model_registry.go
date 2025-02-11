package registry

import (
	"fmt"
	"sync"

	"image/internal/domain/models"
	"image/internal/domain/ports"
	apperrors "image/pkg/errors"
)

// ModelRegistry implements the ModelRegistry interface
type ModelRegistry struct {
	models map[string]models.AIModel
	mu     sync.RWMutex
}

// NewModelRegistry creates a new model registry instance
func NewModelRegistry() ports.ModelRegistry {
	return &ModelRegistry{
		models: make(map[string]models.AIModel),
	}
}

// Register adds a new model to the registry
func (r *ModelRegistry) Register(model models.AIModel) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if model == nil {
		return apperrors.NewInvalidRequestError("Model cannot be nil", nil)
	}

	if _, exists := r.models[model.ID()]; exists {
		return apperrors.NewInvalidRequestError(
			fmt.Sprintf("Model with ID %s already registered", model.ID()),
			nil,
		)
	}

	r.models[model.ID()] = model
	return nil
}

// Get retrieves a model by its ID
func (r *ModelRegistry) Get(id string) (models.AIModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	model, exists := r.models[id]
	if !exists {
		return nil, apperrors.NewInvalidRequestError(
			fmt.Sprintf("Model with ID %s not found", id),
			nil,
		)
	}

	return model, nil
}

// List returns all registered models
func (r *ModelRegistry) List() []models.AIModel {
	r.mu.RLock()
	defer r.mu.RUnlock()

	models := make([]models.AIModel, 0, len(r.models))
	for _, model := range r.models {
		models = append(models, model)
	}

	return models
}

// Validate checks if a model ID is valid
func (r *ModelRegistry) Validate(id string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, exists := r.models[id]; !exists {
		return apperrors.NewInvalidRequestError(
			fmt.Sprintf("Invalid model ID: %s", id),
			nil,
		)
	}

	return nil
}
