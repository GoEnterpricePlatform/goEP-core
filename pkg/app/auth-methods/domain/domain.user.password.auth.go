package domain

import (
	"time"
)

// UserPasswordAuth stores password authentication data.
type UserPasswordAuth struct {
	Password     string     `bson:"-"`
	PasswordHash string     `bson:"password_hash"`
	CreatedAt    *time.Time `bson:"created_at"`
	UpdatedAt    *time.Time `bson:"updated_at"`
}
