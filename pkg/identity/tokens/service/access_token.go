package service

import "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/claim"

// CreateAccessToken generates a signed access token
func (ts *Service) CreateAccessToken(userID string, email string, roles []string, permissions []string) (string, int64, error) {
	claims := claim.NewAccessTokenClaim(userID, email, ts.Issuer, roles, permissions, ts.AccessExpiresIn)
	token, err := claims.GetToken(ts.AccessSecret)
	if err != nil {
		return "", 0, err
	}
	return token, int64(ts.AccessExpiresIn.Seconds()), nil
}

// ParseAccessToken validates and extracts access token claims
func (ts *Service) ParseAccessToken(tokenString string) (*claim.AccessTokenClaims, error) {
	return claim.GetAccessTokenFromJWT(tokenString, ts.AccessSecret)
}
