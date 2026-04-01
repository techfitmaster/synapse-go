package config

import "testing"

func TestValidator_AllPass(t *testing.T) {
	v := NewValidator().
		RequireNonEmpty("JWT_SECRET", "my-real-secret").
		RequireNonEmpty("MYSQL_DSN", "root:pass@tcp(localhost)/db").
		RequireNotDefault("JWT_SECRET", "my-real-secret", "change-me")

	if err := v.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidator_EmptyValue(t *testing.T) {
	v := NewValidator().
		RequireNonEmpty("JWT_SECRET", "").
		RequireNonEmpty("PORT", "8080")

	err := v.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errors) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(ve.Errors), ve.Errors)
	}
}

func TestValidator_DefaultValue(t *testing.T) {
	v := NewValidator().
		RequireNotDefault("JWT_SECRET", "change-me-in-production", "change-me-in-production")

	err := v.Validate()
	if err == nil {
		t.Fatal("expected validation error for default value")
	}
}

func TestValidator_MultipleErrors(t *testing.T) {
	v := NewValidator().
		RequireNonEmpty("JWT_SECRET", "").
		RequireNonEmpty("MYSQL_DSN", "").
		RequireNotDefault("ADMIN_SECRET", "change-me", "change-me")

	err := v.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}

	ve := err.(*ValidationError)
	if len(ve.Errors) != 3 {
		t.Errorf("expected 3 errors, got %d: %v", len(ve.Errors), ve.Errors)
	}
}

func TestValidator_WhitespaceOnly(t *testing.T) {
	v := NewValidator().
		RequireNonEmpty("JWT_SECRET", "   ")

	err := v.Validate()
	if err == nil {
		t.Fatal("expected validation error for whitespace-only value")
	}
}
