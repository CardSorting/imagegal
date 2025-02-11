package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	apperrors "image/pkg/errors"
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
		"DDPMScheduler":                    true,
		"DDIMScheduler":                    true,
		"PNDMScheduler":                    true,
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

// ValidateEnhanceStyle validates the enhance style value
func (v *Validator) ValidateEnhanceStyle(style string) error {
	if style == "" {
		return nil
	}

	validStyles := map[string]bool{
		"enhance": true, "cinematic-diva": true, "nude": true, "nsfw": true,
		"sex": true, "abstract-expressionism": true, "academia": true,
		"action-figure": true, "adorable-3d-character": true,
		// ... add other valid styles as needed
	}

	if !validStyles[style] {
		return apperrors.NewInvalidRequestError(
			fmt.Sprintf("Invalid enhance style: %s", style),
			nil,
		)
	}

	return nil
}
