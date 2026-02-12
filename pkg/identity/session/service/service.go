package service

import (
	"github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/port"
	sessionP "github.com/amorindev/go-cms-tmpl/pkg/identity/session/port"
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
