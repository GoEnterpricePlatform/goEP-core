package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/shared/api/core"
	"github.com/amorindev/go-tmpl/pkg/shared/domain"
)

// RefreshTokenMdw checks and validates refresh tokens from request body
func (m *AuthMiddleware) RefreshTokenMdw(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			core.RespondError(w, domain.NewAppError(domain.ErrCodeInvalidParams, "invalid request body"))
			return
		}

		if req.RefreshToken == "" {
			core.RespondError(w, domain.NewAppError(domain.ErrCodeInvalidParams, "refresh_token is required"))
			return
		}

		c, err := m.TokenSrv.ParseRefreshToken(req.RefreshToken)
		if err != nil {
			core.RespondError(w, domain.ManageError(err, ""))
			return
		}

		ctx := context.WithValue(r.Context(), RefreshTokenClaimsKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
