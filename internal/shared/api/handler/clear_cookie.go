package handler

import "net/http"

// clearCookie removes the cookie by setting an immediate expiration.
func ClearCookie(w http.ResponseWriter, name string, env string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	if env == "prod" {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(w, cookie)
}