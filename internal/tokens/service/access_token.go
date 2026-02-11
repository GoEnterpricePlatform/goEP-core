package service

import "github.com/amorindev/go-cms-tmpl/internal/tokens/claim"

// CreateAccessToken generates a signed access token
func (ts *Service) CreateAccessToken(userID string, email string, roles []string) (string, int64, error) {
	claims := claim.NewAccessTokenClaim(userID, email, ts.Issuer, roles, ts.AccessExpiresIn)
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