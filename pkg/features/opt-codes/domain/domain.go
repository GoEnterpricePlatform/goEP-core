package domain

import "time"

// OtpCodes represents a one-time password (OTP) record
// used for user verification or authentication flows.
// Each code has a specific purpose, an expiration date,
// and a flag indicating whether it has been used.
type OtpCode struct {
	ID        interface{}    `bson:"_id"`
	UserID    interface{}    `bson:"user_id"`
	OptCode   string         `bson:"otp_code"`
	Purpose   OtpCodePurpose `bson:"purpose"`
	Used      bool           `bson:"used"`
	ExpiresAt *time.Time     `bson:"expires_at"`
	CreatedAt *time.Time     `bson:"created_at"`
}

type OtpCodePurpose string

const (
	VerifyEmailPurpose OtpCodePurpose = "verify-email"
)
