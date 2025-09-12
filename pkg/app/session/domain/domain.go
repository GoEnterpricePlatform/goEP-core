package domain

import "time"

// RefreshTokenID: We use this field because when creating a refresh token we need an ID.
// The format will be UUID, independent of the database. We will use this field
// to look up and refresh the token.
//
// RememberMe: Indicates whether to generate a new token with the same expiration
// duration as the original one.
type Session struct {
	ID                interface{} `bson:"_id"`
	UserID            interface{} `bson:"user_id" `
	AccessToken       string      `bson:"-"`
	AccessTokenExpIn  int64       `bson:"-"`
	RefreshTokenID    string      `bson:"refresh_token_id"`
	RefreshToken      string      `bson:"-"`
	RefreshTokenExpIn int64       `bson:"refresh_token_expires_in"`
	Revoked           bool        `bson:"revoked"`
	RememberMe        bool        `bson:"remember_me"`
	CreatedAt         *time.Time  `bson:"create_at"`
}

func NewSession(userID string, rememberMe bool) *Session {
	return &Session{
		UserID:     userID,
		RememberMe: rememberMe,
	}
}
