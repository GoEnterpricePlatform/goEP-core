package handler

import "net/http"

// CreateCookie creates a refresh token cookie with the correct flags
func CreateCookie(refreshToken string, maxAge int, appEnv string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
	}

	if appEnv == "prod" {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}
	return cookie
}
