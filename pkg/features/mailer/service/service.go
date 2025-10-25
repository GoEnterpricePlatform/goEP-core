package service

import "github.com/amorindev/go-tmpl/pkg/features/mailer/port"

var _ port.MailerSrv = &Service{}

type Service struct {
	MailerAdt port.MailerAdt
	AppName   string
}

func NewMailerSrv(mailerAdt port.MailerAdt, appName string) *Service {
	return &Service{
		MailerAdt: mailerAdt,
        AppName: appName,
	}
}

