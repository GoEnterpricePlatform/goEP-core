package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// Validates and ensures the prompt is not empty and within the allowed length.
func validatePrompt(prompt string) error {
	if strings.TrimSpace(prompt) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "prompt is required")
	}
	if len(prompt) > 500 {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "prompt is too long (max 500 characters)")
	}
	return nil
}
