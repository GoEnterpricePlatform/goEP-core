package config

type MailerProvider string
type DBProvider string

const (
	MailGmail  MailerProvider = "gmail"
	MailResend MailerProvider = "resend"
)

const (
	DBMongo DBProvider = "mongo"
)

type AppStack struct {
	DB   DBProvider
	Mail MailerProvider
}

func NewAppStack(db DBProvider, mail MailerProvider) *AppStack {
	return &AppStack{
		DB:   db,
		Mail: mail,
	}
}