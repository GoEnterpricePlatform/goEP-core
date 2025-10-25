package port

type MailerAdt interface {
	Send(to, subject, htmlBody string) error
}

type MailerSrv interface {
	SendVerifyEmail(email string, code string) error
}
