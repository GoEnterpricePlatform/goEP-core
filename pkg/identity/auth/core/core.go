package core

import (
	sessionC "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/core"
	sessionD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/domain"
	userC "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/core"
	userD "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
)

// AuthResp represents the unified authentication response structure
// This structure handles all authentication scenarios (sign up, sign in, verification, etc.)
// Fields are conditionally populated based on the authentication flow and user state
type AuthResp struct {
	// Session contains authentication tokens (null for sign up before email verification)
	Session *sessionC.SessionCore `json:"session"`

	// User contains user profile information (always present)
	User *userC.UserCore `json:"user"`

	// OtpID represents the ID of the OTP (One-Time Password) used for email verification.
	// It is only present during sign-up before verification is completed.
	OtpID *string `json:"otp_id"`
}

// NewAuthResp creates a new AuthResp with flexible field population
// This function handles all authentication scenarios
func NewAuthResp(user *userD.User, session *sessionD.Session, otpID string) *AuthResp {
	resp := &AuthResp{}

	if user != nil {
		resp.User = &userC.UserCore{
			ID:            user.ID,
			Email:         user.Email,
			EmailVerified: user.EmailVerified,
			ImgUrl:        user.ImgUrl,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		}
	}

	if session != nil {
		resp.Session = &sessionC.SessionCore{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			ExpiresIn:    session.AccessTokenExpIn,
		}
	}

	if otpID != "" {
		resp.OtpID = &otpID
	}

	return resp
}

// NewSignUpResp creates response for user registration
func NewSignUpResp(user *userD.User, otpID string) *AuthResp {
	return NewAuthResp(user, nil, otpID)
}

// NewSignInResp creates response for successful sign in
func NewSignInResp(user *userD.User, session *sessionD.Session, otpID string) *AuthResp {
	return NewAuthResp(user, session, otpID)
}

// NewRefreshTokenResp creates response for successful refresh token
func NewRefreshTokenResp(session *sessionD.Session) *AuthResp {
	return NewAuthResp(nil, session, "")
}
