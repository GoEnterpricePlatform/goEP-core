package model

import (
	"github.com/google/uuid"
	"time"
)

// UserSQLModel matches tb_identity_user. The UUID is a persistence detail;
// the domain continues to expose IDs as strings across database providers.
type UserSQLModel struct {
	ID            uuid.UUID
	Email         string
	EmailVerified bool
	IsActive      bool
	ImgPath       *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// UserPasswordSQLModel maps the optional one-to-one tb_identity_pass_auth row.
type UserPasswordSQLModel struct {
	UserID       uuid.UUID
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRoleSQLModel maps the user-role join table used by SQL persistence.
type UserRoleSQLModel struct {
	UserID uuid.UUID
	RoleID uuid.UUID
}
