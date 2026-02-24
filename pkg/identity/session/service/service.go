package service

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/port"
	sessionP "github.com/GoEnterpricePlatform/goEP-core/pkg/identity/session/port"
)

var _ sessionP.SessionSrv = &Service{}

type Service struct {
	SessionRepo sessionP.SessionRepo
	TokenSrv    port.TokenSrv
}

func NewSessionSrv(sessionRepo sessionP.SessionRepo, tokenSrv port.TokenSrv) *Service {
	return &Service{
		SessionRepo: sessionRepo,
		TokenSrv:    tokenSrv,
	}
}
