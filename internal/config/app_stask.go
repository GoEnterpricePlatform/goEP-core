package config

type MailerProvider string
type DBProvider string
type FSProvider string

const (
	MailGmail  MailerProvider = "gmail"
	MailResend MailerProvider = "resend"
)

const (
	DBMongo DBProvider = "mongo"
)

const (
	FSMinio FSProvider = "minio"
)

type AppStack struct {
	DB      DBProvider
	Mail    MailerProvider
	FileStg FSProvider
}

func NewAppStack(db DBProvider, mail MailerProvider, fileStg FSProvider) *AppStack {
	return &AppStack{
		DB:   db,
		Mail: mail,
		FileStg: fileStg,
	}
}
