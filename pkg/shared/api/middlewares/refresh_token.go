package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
)

// RefreshTokenMdw checks and validates refresh tokens from request body
// NOTE: This middleware is no longer used because refresh tokens are now
// retrieved directly from cookies and validated inside the corresponding
// handlers (e.g., /auth/refresh and /auth/logout). Since only these two
// endpoints require refresh token validation, it is simpler and clearer to
// handle this logic directly in each handler rather than maintaining a
// dedicated middleware. This code is kept for reference only and can be
// safely removed later if not needed.
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
