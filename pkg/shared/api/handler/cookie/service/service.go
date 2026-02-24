package service

import (
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler/cookie/port"
)

var _ port.CookieSrv = &Service{}

type Service struct {
	appEnv             string
	JwtAccessCookieDur time.Duration
}

func NewCookieSrv(appEnv string, jwtAccessCookieDur time.Duration) *Service {
	return &Service{
		appEnv:             appEnv,
		JwtAccessCookieDur: jwtAccessCookieDur,
	}
}
