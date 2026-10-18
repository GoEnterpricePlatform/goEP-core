package port

import (
	"net/http"
	"time"
)

type CookieSrv interface {
	SetAccessToken(w http.ResponseWriter, token string)
	SetRefreshToken(w http.ResponseWriter, token string, dur time.Duration)

	ClearForTemplate(w http.ResponseWriter, name string)
}
