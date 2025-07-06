package password

import (
	"fmt"

	"mira/common/xerrors"

	"golang.org/x/crypto/bcrypt"
)

// Generate password
func Generate(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash: %w", err)
	}

	return string(hash), nil
}

// Verify password
func Verify(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return xerrors.ErrMismatchedPassword
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}
	return nil
}
