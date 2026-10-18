package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
)

// ValidateCreate checks the fields required for creating a post
func validatePostFields(name string) error {
	if strings.TrimSpace(name) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "name is required")
	}
	if len(name) > 100 {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "name cannot exceed 100 characters")
	}
	return nil
}
