package gmailsmtp

import (
	"fmt"
	"net/smtp"
)

func (a *Adapter) Send(to, subject, htmlBody string) error {
	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s",
		a.From,
		to,
		subject,
		htmlBody,
	)
	return smtp.SendMail(a.Addr, a.Client, a.From, []string{to}, []byte(message))
}
