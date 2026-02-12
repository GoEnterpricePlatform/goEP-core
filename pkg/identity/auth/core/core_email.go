package core

import (
	"strings"

	dShared "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"github.com/go-playground/validator/v10"
)

type EmailReq struct {
	Email string `json:"email"`
}

func (req *EmailReq) IsEmailValid() error {
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
	return nil
}
