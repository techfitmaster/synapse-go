package config

import (
	"fmt"
	"strings"
)

// ValidationError collects multiple configuration validation failures.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config validation failed:\n  - %s", strings.Join(e.Errors, "\n  - "))
}

// Validator checks configuration values at startup and collects errors.
type Validator struct {
	errors []string
}

// NewValidator creates a configuration validator.
func NewValidator() *Validator {
	return &Validator{}
}

// RequireNonEmpty checks that the given value is not empty. Adds an error if it is.
func (v *Validator) RequireNonEmpty(name, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, fmt.Sprintf("%s is required but empty", name))
	}
	return v
}

// RequireNotDefault checks that the given value differs from a known insecure default.
func (v *Validator) RequireNotDefault(name, value, defaultVal string) *Validator {
	if value == defaultVal {
		v.errors = append(v.errors, fmt.Sprintf("%s must be changed from default value", name))
	}
	return v
}

// Validate returns a ValidationError if any checks failed, or nil if all passed.
func (v *Validator) Validate() error {
	if len(v.errors) == 0 {
		return nil
	}
	return &ValidationError{Errors: v.errors}
}
