package gmailsmtp

import (
	"net/smtp"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/mailer/port"
)

var _ port.MailerAdt = &Adapter{}

type Adapter struct {
	Client smtp.Auth
	Addr   string
	From   string
}

func NewGmailSmtpAdt(client smtp.Auth, addr string, from string) *Adapter {
	return &Adapter{
		Client: client,
		Addr:   addr,
		From:   from,
	}
}
