package service

import (
	"time"

	"github.com/amorindev/go-cms-tmpl/internal/tokens/port"
)

var _ port.TokenSrv = &Service{}

type Service struct {
	AccessSecret             string
	RefreshSecret            string
	AccessExpiresIn          time.Duration
	RefreshExpiresIn         time.Duration
	RefreshExpiresInRemember time.Duration
	Issuer                   string
}

func NewTokenSrv(accessSecret string, refreshSecret string, accessExpiresIn time.Duration, refreshExpiresIn time.Duration, refreshExpiresInRemember time.Duration, issuer string) *Service {
	return &Service{
		AccessSecret:             accessSecret,
		RefreshSecret:            refreshSecret,
		AccessExpiresIn:          accessExpiresIn,
		RefreshExpiresIn:         refreshExpiresIn,
		Issuer:                   issuer,
		RefreshExpiresInRemember: refreshExpiresInRemember,
	}
}
