package service

import "github.com/amorindev/go-cms-tmpl/pkg/identity/mailer/service/templates"

type VerifyEmailData struct {
	Name    string
	Subject string
	Code    string
}

func (s *Service) SendVerifyEmail(email string, code string) error {

	data := VerifyEmailData{
		Name:    s.AppName,
		Subject: email,
		Code:    code,
	}

	tmplString, err := templates.LoadTemplate("pkg/identity/mailer/service/templates/verify-email.html", data)
	if err != nil {
		return err
	}

	subject := "Verify Your Email Address"

	return s.MailerAdt.Send(email, subject, tmplString)
}
