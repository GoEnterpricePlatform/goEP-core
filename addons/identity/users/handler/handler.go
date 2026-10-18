package handler

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/users/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares"
)

type Handler struct {
	UserSrv    port.UserSrv
	AuthApiMdw *middlewares.AuthMiddleware
}

func NewUserHandler(server *http.ServeMux, userSrv port.UserSrv, authApiMdw *middlewares.AuthMiddleware) *Handler {
	h := &Handler{
		UserSrv:    userSrv,
		AuthApiMdw: authApiMdw,
	}

	server.Handle("POST /users/{userId}/avatar", h.AuthApiMdw.AccessTokenMdw(h.UploadAvatar))

	return h
}
