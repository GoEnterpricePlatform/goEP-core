package domain

import "time"

// RefreshTokenID: We use this field because when creating a refresh token we need an ID.
// The format will be UUID, independent of the database. We will use this field
// to look up and refresh the token.
//
// RememberMe: Indicates whether to generate a new token with the same expiration
// duration as the original one.
type Session struct {
	ID                string
	UserID            string
	AccessToken       string
	AccessTokenExpIn  int64
	RefreshTokenID    string
	RefreshToken      string
	RefreshTokenExpIn int64
	Revoked           bool
	RememberMe        bool
	CreatedAt         *time.Time
}

func NewSession(userID string, rememberMe bool) *Session {
	return &Session{
		UserID:     userID,
		RememberMe: rememberMe,
	}
}
