package core

import (
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/features/users/domain"
)

// UserCore represents the user data returned to the client, omitting sensitive fields like passwords.
type UserCore struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	ImgUrl        *string    `json:"img_url"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

// NewFromUserDomain converts a domain.User object to a UserCore for API responses.
func NewFromUserDomain(user *domain.User) UserCore {
	return UserCore{
		ID:            user.ID.(string),
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}
