package service

import "github.com/amorindev/go-tmpl/pkg/shared/api/handler/cookie/port"

var _ port.CookieSrv = &Service{}

type Service struct {
	appEnv string
}

func NewCookieSrv(appEnv string) *Service {
	return &Service{appEnv: appEnv}
}
