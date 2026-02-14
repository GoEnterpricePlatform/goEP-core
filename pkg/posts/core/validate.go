package core

import (
	"strings"

	"github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

// ValidateCreate checks the fields required for creating a post
func validatePostFields(title string, desc *string) error {
	if strings.TrimSpace(title) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "title is required")
	}
	if len(title) > 100 {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "title cannot exceed 100 characters")
	}
	if desc != nil && len(*desc) > 255 {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "desc cannot exceed 255 characters")
	}
	return nil
}
