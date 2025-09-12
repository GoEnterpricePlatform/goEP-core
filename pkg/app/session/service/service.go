package service

import (
	"github.com/amorindev/go-tmpl/internal/auth"
	"github.com/amorindev/go-tmpl/pkg/app/session/port"
)

var _ port.SessionSrv = &Service{}

type Service struct {
	SessionRepo port.SessionRepo
	TokenSrv    *auth.TokenSrv
}

func NewSessionSrv(sessionRepo port.SessionRepo, tokenSrv *auth.TokenSrv) *Service {
	return &Service{
		SessionRepo: sessionRepo,
		TokenSrv:    tokenSrv,
	}
}
