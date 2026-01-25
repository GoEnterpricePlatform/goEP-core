package gmailsmtp

import (
	"net/smtp"
)

func NewGmailSmtpClient(username string, password string, host string) smtp.Auth {
	return smtp.PlainAuth(
		"",
		username,
		password,
		host,
	)
}
