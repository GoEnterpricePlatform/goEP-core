package service

import "net/http"

func (s *Service) ClearForTemplate(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   s.appEnv == "prod",
	}
	http.SetCookie(w, cookie)
}
