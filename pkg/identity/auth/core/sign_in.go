package core

import (
	"strings"

	dShared "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"github.com/go-playground/validator/v10"
)

// SignInReq represents the sign in request structure
type SignInReq struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

// IsSignInValid performs validation on the sign in request
// Returns an error if any validation fails
func (req SignInReq) IsSignInValid() error {
	// Validate email field is not empty
	if strings.TrimSpace(req.Email) == "" {
		return dShared.NewAppError(dShared.ErrCodeInvalidParams, "email is required")
	}

	// Validate email format using validator
	validate := validator.New()
	err := validate.Var(req.Email, "email")
	if err != nil {
		return dShared.NewAppError(dShared.ErrCodeInvalidParams, "invalid email format")
	}

	// Validate password field is not empty
	if strings.TrimSpace(req.Password) == "" {
		return dShared.NewAppError(dShared.ErrCodeInvalidParams, "password is required")
	}

	return nil
}
