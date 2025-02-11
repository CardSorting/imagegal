package validation

import (
	"fmt"
	"reflect"
	"strings"

	apperrors "image/pkg/errors"

	"github.com/go-playground/validator/v10"
)

// Validator represents a validator instance
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Register custom validation tags if needed
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validate: v,
	}
}

// Validate validates the provided struct
func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return v.handleValidationErrors(validationErrors)
		}
		return apperrors.NewInvalidRequestError("Invalid request", err)
	}
	return nil
}

// handleValidationErrors processes validation errors into a user-friendly format
func (v *Validator) handleValidationErrors(errors validator.ValidationErrors) error {
	var errorMessages []string

	for _, err := range errors {
		field := err.Field()
		switch err.Tag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("%s is required", field))
		case "min":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param()))
		case "max":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be less than or equal to %s", field, err.Param()))
		case "oneof":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be one of [%s]", field, err.Param()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s failed validation: %s", field, err.Tag()))
		}
	}

	return apperrors.NewInvalidRequestError(
		fmt.Sprintf("Validation failed: %s", strings.Join(errorMessages, "; ")),
		nil,
	)
}

// ValidateScheduler validates the scheduler name
func (v *Validator) ValidateScheduler(scheduler string) error {
	validSchedulers := map[string]bool{
		"DDPMScheduler":                   true,
		"DDIMScheduler":                   true,
		"PNDMScheduler":                   true,
		"LMSDiscreteScheduler":            true,
		"EulerDiscreteScheduler":          true,
		"EulerAncestralDiscreteScheduler": true,
		"DPMSolverMultistepScheduler":     true,
		"HeunDiscreteScheduler":           true,
		"KDPM2DiscreteScheduler":          true,
		"DPMSolverSinglestepScheduler":    true,
		"KDPM2AncestralDiscreteScheduler": true,
		"UniPCMultistepScheduler":         true,
		"DDIMInverseScheduler":            true,
		"DEISMultistepScheduler":          true,
		"IPNDMScheduler":                  true,
		"KarrasVeScheduler":               true,
		"ScoreSdeVeScheduler":             true,
		"LCMScheduler":                    true,
	}

	if scheduler != "" && !validSchedulers[scheduler] {
		return apperrors.NewInvalidRequestError(
			fmt.Sprintf("Invalid scheduler: %s", scheduler),
			nil,
		)
	}

	return nil
}

// ValidateEnhancePrompt validates the enhance prompt value
func (v *Validator) ValidateEnhancePrompt(value string) error {
	if value == "" {
		return nil
	}

	validValues := map[string]bool{
		"yes": true,
		"no":  true,
	}

	if !validValues[value] {
		return apperrors.NewInvalidRequestError(
			fmt.Sprintf("Invalid enhance_prompt value: %s (must be 'yes' or 'no')", value),
			nil,
		)
	}

	return nil
}
