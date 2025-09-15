package handler

import (
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/app/users/port"
)

type Handler struct {
	UserSrv port.UserSrv
}

func NewUserHandler(server *http.ServeMux, userSrv port.UserSrv) *Handler {
	h := &Handler{
		UserSrv: userSrv,
	}

	server.HandleFunc("POST /users/{userId}/avatar", h.UploadAvatar)

	return h
}
