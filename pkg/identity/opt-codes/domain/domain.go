package domain

import "time"

// OtpCodes represents a one-time password (OTP) record
// used for user verification or authentication flows.
// Each code has a specific purpose, an expiration date,
// and a flag indicating whether it has been used.
type OtpCode struct {
	ID        string
	UserID    string
	OptCode   string
	Purpose   OtpCodePurpose
	Used      bool
	ExpiresAt *time.Time
	CreatedAt *time.Time
}

type OtpCodePurpose string

const (
	VerifyEmailPurpose OtpCodePurpose = "verify-email"
)
