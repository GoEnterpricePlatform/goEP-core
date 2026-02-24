package claim

import (
	"errors"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// RefreshTokenClaims defines the payload for the refresh token
type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// NewRefreshTokenClaim creates a new set of claims for a refresh token
func NewRefreshTokenClaim(userID string, expiresIn time.Duration) *RefreshTokenClaims {
	return &RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
	}
}

// GetToken generates and signs a JWT refresh token
func (c *RefreshTokenClaims) GetToken(refreshSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return signed, err
}

// GetRefreshTokenFromJWT validates and extracts claims from a refresh token
func GetRefreshTokenFromJWT(tokenString, refreshSecret string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
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

		return nil, fmt.Errorf("invalid refresh token: %s", err.Error())
	}

	if !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	claim, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, domain.ErrTokenInvalidClaim
	}

	return claim, nil
}
