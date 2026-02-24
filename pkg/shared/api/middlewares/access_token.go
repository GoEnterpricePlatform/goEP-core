package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
)

// AccessTokenMdw checks and validates access tokens from Authorization header
func (m *AuthMiddleware) AccessTokenMdw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString, err := tokenFromAuthorization(authHeader)

		if err != nil {
			core.RespondError(w, domain.ManageError(err, ""))
			return
		}

		c, err := m.TokenSrv.ParseAccessToken(tokenString)
		if err != nil {
			core.RespondError(w, domain.ManageError(err, ""))
			return
		}

		ctx := context.WithValue(r.Context(), AccessTokenClaimsIDKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// tokenFromAuthorization extracts the token from "Bearer <token>"
func tokenFromAuthorization(authorization string) (string, error) {
	if authorization == "" {
		return "", domain.ErrAuthHeaderMissing
	}

	if !strings.HasPrefix(authorization, "Bearer") {
		return "", domain.ErrAuthHeaderInvalid
	}

	l := strings.Split(authorization, " ")
	if len(l) != 2 {
		return "", domain.ErrAuthHeaderInvalid
	}

	return l[1], nil
}
