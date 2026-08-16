package service

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/port"
)

var _ port.UserSrv = &Service{}

type Service struct {
	UserRepo    port.UserRepo
	UserFileStg port.UserFileStg
	MdlName     string
}

func NewUserSrv(userRepo port.UserRepo, userFileStg port.UserFileStg, mdlName string) *Service {
	return &Service{
		UserRepo:    userRepo,
		UserFileStg: userFileStg,
		MdlName: mdlName,
	}
}
