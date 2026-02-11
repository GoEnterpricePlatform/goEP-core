package service

import (
	"github.com/amorindev/go-cms-tmpl/internal/tokens/port"
	sessionP "github.com/amorindev/go-cms-tmpl/pkg/features/session/port"
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
