package crypto

import (
	"testing"
)

func TestGenerateNumericCode_Length(t *testing.T) {
	code, err := GenerateNumericCode(6)
	if err != nil {
		t.Fatalf("GenerateNumericCode(6) error: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("len = %d, want 6", len(code))
	}
}

func TestGenerateNumericCode_AllDigits(t *testing.T) {
	code, err := GenerateNumericCode(6)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("non-digit character %q in code %q", c, code)
		}
	}
}

func TestGenerateNumericCode_Randomness(t *testing.T) {
	code1, _ := GenerateNumericCode(6)
	code2, _ := GenerateNumericCode(6)
	// Extremely unlikely to be equal (1 in 1,000,000)
	if code1 == code2 {
		t.Logf("warning: two codes are equal (%s), extremely unlikely but not impossible", code1)
	}
}

func TestGenerateNumericCode_PaddedWithZeros(t *testing.T) {
	// Generate many codes to check padding works for small numbers
	for i := 0; i < 100; i++ {
		code, err := GenerateNumericCode(6)
		if err != nil {
			t.Fatal(err)
		}
		if len(code) != 6 {
			t.Errorf("code %q has length %d, want 6", code, len(code))
		}
	}
}

func TestGenerateNumericCode_DifferentLengths(t *testing.T) {
	tests := []int{4, 6, 8}
	for _, length := range tests {
		code, err := GenerateNumericCode(length)
		if err != nil {
			t.Fatalf("GenerateNumericCode(%d) error: %v", length, err)
		}
		if len(code) != length {
			t.Errorf("GenerateNumericCode(%d) len = %d", length, len(code))
		}
	}
}
