package model

import (
	"github.com/google/uuid"
	"time"
)

// SessionSQLModel maps tb_identity_session. UserID is a relational foreign
// key and the token secret remains outside persistent session state.
type SessionSQLModel struct {
	ID                    uuid.UUID
	UserID                uuid.UUID
	RefreshTokenID        string
	RefreshTokenExpiresIn int64
	Revoked               bool
	RememberMe            bool
	CreatedAt             time.Time
}
