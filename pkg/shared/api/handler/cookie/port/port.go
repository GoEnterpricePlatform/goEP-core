package port

import "net/http"

type CookieSrv interface {
	CreateForTemplate(name, value string, maxAge int) *http.Cookie
	ClearForTemplate(w http.ResponseWriter, name string)
}
