package resend

import (
	"github.com/amorindev/go-cms-tmpl/pkg/features/mailer/port"
	"github.com/resend/resend-go/v2"
)

var _ port.MailerAdt = &Adapter{}

type Adapter struct {
	Client *resend.Client
	From   string
}

func NewResendAdt(client *resend.Client, from string) *Adapter {
	return &Adapter{
		Client: client,
		From:   from,
	}
}
