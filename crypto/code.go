package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateNumericCode generates a cryptographically random numeric code of the given length.
// For example, GenerateNumericCode(6) returns a string like "042817".
func GenerateNumericCode(length int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("generate code: %w", err)
	}
	return fmt.Sprintf("%0*d", length, n), nil
}
