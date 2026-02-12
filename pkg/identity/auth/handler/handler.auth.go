package handler

import (
	"net/http"

	tokenP "github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/port"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/auth/port"
)

type Handler struct {
	AuthSrv  port.AuthSrv
	TokenSrv tokenP.TokenSrv
	AppEnv   string
}

func NewAuthHandler(server *http.ServeMux, authSrv port.AuthSrv, tokenSrv tokenP.TokenSrv, appEnv string) *Handler {
	h := &Handler{
		AuthSrv:  authSrv,
		AppEnv:   appEnv,
		TokenSrv: tokenSrv,
	}

	server.HandleFunc("POST /auth/sign-up", h.SignUp)
	server.HandleFunc("POST /auth/sign-in", h.SignIn)
	server.HandleFunc("POST /auth/resend-verify-email", h.ResendVerifyEmail)
	server.HandleFunc("POST /auth/verify-email", h.VerifyEmail)
	server.HandleFunc("POST /auth/sign-out", h.SignOut)
	server.HandleFunc("POST /auth/refresh-token", h.RefreshToken)
	server.HandleFunc("GET /auth/get-session", h.GetSession)

	return h
}
