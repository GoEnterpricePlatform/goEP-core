package encryption

import (
	"errors"
	"fmt"

	"github.com/amorindev/go-tmpl/pkg/shared/domain"
	"golang.org/x/crypto/bcrypt"
)

// CheckPassword compares a plain-text password with its hashed value.
// Returns nil if they match, or an error if they do not match or if
// verification fails.
func CheckPassword(password string, hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("%w: error comparing hash: %w", domain.ErrPassDoNotMatch, err)
		}
		return fmt.Errorf("error comparing hash: %w", err)
	}
	return nil
}
