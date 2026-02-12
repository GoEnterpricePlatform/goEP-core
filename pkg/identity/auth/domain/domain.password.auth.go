package domain

import (
	"time"
)

// UserPasswordAuth stores password authentication data.
type UserPasswordAuth struct {
	Password     string
	PasswordHash string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}
