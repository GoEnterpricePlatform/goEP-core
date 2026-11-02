package model

import (
	"github.com/google/uuid"
	"time"
)

// OtpCodeSQLModel maps tb_identity_otp_code. UserID is the SQL foreign key
// required by this relational store and does not change the domain contract.
type OtpCodeSQLModel struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	OtpCode   string
	Purpose   string
	Used      bool
	ExpiresAt time.Time
	CreatedAt time.Time
}
