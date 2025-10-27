package handler

import (
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/features/auth/port"
)

type Handler struct {
	AuthSrv port.AuthSrv
	AppEnv  string
}

func NewAuthHandler(server *http.ServeMux, authSrv port.AuthSrv, appEnv string) *Handler {
	h := &Handler{
		AuthSrv: authSrv,
		AppEnv:  appEnv,
	}

	server.HandleFunc("POST /auth/sign-up", h.SignUp)
	server.HandleFunc("POST /auth/sign-in", h.SignIn)
	server.HandleFunc("POST /auth/resend-verify-email", h.ResendVerifyEmail)
	server.HandleFunc("POST /auth/verify-email", h.VerifyEmail)

	return h
}
