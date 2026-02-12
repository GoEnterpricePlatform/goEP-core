package service

import "github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/claim"

// CreateRefreshToken generates a signed refresh token
func (ts *Service) CreateRefreshToken(userID string, rememberMe bool) (string, string, int64, error) {
	ttl := ts.RefreshExpiresIn
	if rememberMe {
		ttl = ts.RefreshExpiresInRemember
	}

	claims := claim.NewRefreshTokenClaim(userID, ttl)

	token, err := claims.GetToken(ts.RefreshSecret)
	if err != nil {
		return "", "", 0, err
	}
	return claims.ID, token, int64(ttl.Seconds()), nil
}

// ParseRefreshToken validates and extracts refresh token claims
func (ts *Service) ParseRefreshToken(tokenString string) (*claim.RefreshTokenClaims, error) {
	return claim.GetRefreshTokenFromJWT(tokenString, ts.RefreshSecret)
}