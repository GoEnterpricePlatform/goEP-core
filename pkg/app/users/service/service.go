package service

import (
	"github.com/amorindev/go-tmpl/pkg/app/users/port"
)

var _ port.UserSrv = &Service{}

type Service struct {
	UserRepo   port.UserRepo
	UserFileStg port.UserFileStg
}

func NewUserSrv(userRepo port.UserRepo, userFileStg port.UserFileStg) *Service {
	return &Service{
		UserRepo:   userRepo,
		UserFileStg: userFileStg,
	}
}
