package crypto

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("Abcd1234")
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}
	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}
	if hash == "Abcd1234" {
		t.Error("HashPassword() returned plaintext")
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	hash, _ := HashPassword("Abcd1234")
	if err := CheckPassword(hash, "Abcd1234"); err != nil {
		t.Errorf("CheckPassword() should pass for correct password: %v", err)
	}
}

func TestCheckPassword_Wrong(t *testing.T) {
	hash, _ := HashPassword("Abcd1234")
	if err := CheckPassword(hash, "WrongPass1"); err == nil {
		t.Error("CheckPassword() should fail for wrong password")
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name    string
		pw      string
		wantErr bool
	}{
		{"valid", "Abcd1234", false},
		{"valid_long", "MySecurePass123!", false},
		{"too_short", "Ab1", true},
		{"no_letter", "12345678", true},
		{"no_digit", "Abcdefgh", true},
		{"exactly_8_valid", "Abcdefg1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.pw)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordStrength(%q) error = %v, wantErr %v", tt.pw, err, tt.wantErr)
			}
		})
	}
}
