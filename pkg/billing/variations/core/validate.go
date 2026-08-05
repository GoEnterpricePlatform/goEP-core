package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// ValidateCreate checks the fields required for creating a varOption
func validateVarOptionFields(label string, value *string) error {
	if strings.TrimSpace(label) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "label is required")
	}
	if value != nil && strings.TrimSpace(*value) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "value is required")
	}
	return nil
}