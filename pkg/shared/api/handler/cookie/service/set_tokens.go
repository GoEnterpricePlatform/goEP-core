package service

import (
	"net/http"
	"time"
)

// SetAccessToken sets the access token cookie.
// The cookie lifetime is intentionally longer than the JWT expiration
// to allow the backend to detect an expired access token and perform
// automatic refresh before the browser removes the cookie.
func (s *Service) SetAccessToken(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(s.JwtAccessCookieDur.Seconds()),
		SameSite: http.SameSiteLaxMode,
		Secure:   s.appEnv == "prod",
	}

	http.SetCookie(w, cookie)
}

// SetRefreshToken sets the refresh token cookie.
// The duration (dur) must be passed dynamically because the refresh
// TTL depends on the login context (normal vs remember me).
// The cookie must live exactly as long as the refresh JWT generated,
// otherwise the browser could delete it before the token actually expires.
func (s *Service) SetRefreshToken(w http.ResponseWriter, token string, dur time.Duration) {
	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(dur.Seconds()),
		SameSite: http.SameSiteLaxMode,
		Secure:   s.appEnv == "prod",
	}

	http.SetCookie(w, cookie)
}
