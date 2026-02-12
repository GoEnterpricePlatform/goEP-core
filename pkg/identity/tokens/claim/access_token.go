package claim

import (
	"errors"
	"fmt"
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"github.com/golang-jwt/jwt/v5"
)

// AccessTokenClaims defines the payload for the access token
type AccessTokenClaims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Roles       []string `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

// NewAccessTokenClaim creates a new set of claims for an access token
func NewAccessTokenClaim(userID string, email string, issuer string, roles []string, permissions []string, expiresIn time.Duration) *AccessTokenClaims {
	return &AccessTokenClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
	}
}

// GetToken generates and signs a JWT access token
func (c *AccessTokenClaims) GetToken(accessSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := token.SignedString([]byte(accessSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return signed, nil
}

// GetAccessTokenFromJWT validates and extracts claims from a JWT access token
func GetAccessTokenFromJWT(tokenString, accessSecret string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrTokenExpired
		}

		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, domain.ErrTokenSignature
		}

		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, domain.ErrTokenMalformed
		}

		return nil, fmt.Errorf("invalid access token: %s", err.Error())
	}

	if !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	claim, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, domain.ErrTokenInvalidClaim
	}

	return claim, nil
}
