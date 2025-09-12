package core

import (
	sessionC "github.com/amorindev/go-tmpl/pkg/app/session/core"
	sessionD "github.com/amorindev/go-tmpl/pkg/app/session/domain"
	userC "github.com/amorindev/go-tmpl/pkg/app/users/core"
	userD "github.com/amorindev/go-tmpl/pkg/app/users/domain"
)

// AuthResp represents the unified authentication response structure
// This structure handles all authentication scenarios (sign up, sign in, verification, etc.)
// Fields are conditionally populated based on the authentication flow and user state
type AuthResp struct {
	// Session contains authentication tokens (null for sign up before email verification)
	Session *sessionC.SessionCore `json:"session"`

	// User contains user profile information (always present)
	User *userC.UserCore `json:"user"`
}

// NewAuthResp creates a new AuthResp with flexible field population
// This function handles all authentication scenarios
func NewAuthResp(user *userD.User, session *sessionD.Session) *AuthResp {
	resp := &AuthResp{
		User: &userC.UserCore{
			ID:            user.ID.(string),
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		},
	}

	if session != nil {
		resp.Session = &sessionC.SessionCore{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			ExpiresIn:    session.RefreshTokenExpIn,
		}
	}
	return resp
}

// NewSignUpResp creates response for user registration
func NewSignUpResp(user *userD.User) *AuthResp {
	return NewAuthResp(user, nil)
}

// NewSignInResp creates response for successful sign in
func NewSignInResp(user *userD.User, session *sessionD.Session) *AuthResp {
	return NewAuthResp(user, session)
}
