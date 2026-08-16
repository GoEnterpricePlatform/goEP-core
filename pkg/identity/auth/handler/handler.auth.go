package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/port"
	tokenP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/port"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"
)

type Handler struct {
	AuthSrv    port.AuthSrv
	TokenSrv   tokenP.TokenSrv
	AppEnv     string
	AuthApiMdw *middlewares.AuthMiddleware
}

func NewAuthHandler(
	server *http.ServeMux,
	authSrv port.AuthSrv,
	tokenSrv tokenP.TokenSrv,
	appEnv string,
	authApiMdw *middlewares.AuthMiddleware,
) *Handler {
	h := &Handler{
		AuthSrv:    authSrv,
		AppEnv:     appEnv,
		TokenSrv:   tokenSrv,
		AuthApiMdw: authApiMdw,
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
