package validate

import "testing"

func TestEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"user@example.com", false},
		{"a.b+c@domain.co", false},
		{"test@sub.domain.com", false},
		{"", true},
		{"@domain.com", true},
		{"user@", true},
		{"user@.com", true},
		{"user domain.com", true},
		{"user@domain", true},
	}
	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			err := Email(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestPhone(t *testing.T) {
	tests := []struct {
		phone   string
		wantErr bool
	}{
		{"+8613800138000", false},
		{"+14155551234", false},
		{"8613800138000", false},
		{"1234567", false},
		{"", true},
		{"123", true},
		{"+0123456", true},
		{"abc", true},
	}
	for _, tt := range tests {
		t.Run(tt.phone, func(t *testing.T) {
			err := Phone(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Phone(%q) error = %v, wantErr %v", tt.phone, err, tt.wantErr)
			}
		})
	}
}
