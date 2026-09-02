package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// /v1/admin/sign-in was updated to /v1/admin/auth/sign-in
// therefore if you are using web/standard-web it will be re-established 
// you should update the paths to /v1/admin/auth/sign-in
// Why will it be the current form?
func (m *MdwSrvTmpl) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to read the access token from cookie
		accessCookie, err := r.Cookie("accessToken")
		if err != nil {
			m.CookieSrv.ClearForTemplate(w, "accessToken")
			http.Redirect(w, r, "/v1/admin/auth/sign-in", http.StatusFound)
			return
		}

		// Validate and parse the access token
		c, err := m.TokenSrv.ParseAccessToken(accessCookie.Value)
		if err != nil {
			if errors.Is(err, domain.ErrTokenExpired) {
				// Try to read the refresh token from cookie
				refreshCookie, err := r.Cookie("refreshToken")
				if err != nil {
					m.CookieSrv.ClearForTemplate(w, "refreshToken")
					http.Redirect(w, r, "/v1/admin/auth/sign-in", http.StatusFound)
					return
				}

				refreshClaims, err := m.TokenSrv.ParseRefreshToken(refreshCookie.Value)
				if err != nil {
					http.Redirect(w, r, "/v1/admin/auth/sign-in", http.StatusFound)
					return
				}

				session, err := m.AuthSrv.RefreshToken(r.Context(), refreshClaims.ID, refreshClaims.UserID)
				if err != nil {
					http.Redirect(w, r, "/v1/admin/auth/sign-in", http.StatusFound)
					return
				}

				m.CookieSrv.SetAccessToken(w, session.AccessToken)

				claims, err := m.TokenSrv.ParseAccessToken(session.AccessToken)
				if err != nil {
					http.Redirect(w, r, "/v1/admin/auth/sign-in", http.StatusFound)
					return
				}

				ctx := context.WithValue(r.Context(), AccessTokenClaimsTmplIDKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			m.CookieSrv.ClearForTemplate(w, "accessToken")
			http.Redirect(w, r, "/v1/admin/auth/sign-in", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), AccessTokenClaimsTmplIDKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
