package domain

import (
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/features/auth/domain"
)

// User represents a user in the system
type User struct {
	ID            string
	Email         string
	EmailVerified bool
	IsActive      bool
	ImgUrl        *string
	ImgPath       *string
	UserPassAuth  *domain.UserPasswordAuth
	RoleIDs       []string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

func NewUser(email string, password string) *User {
	return &User{
		Email: email,
		UserPassAuth: &domain.UserPasswordAuth{
			Password: password,
		},
	}
}
