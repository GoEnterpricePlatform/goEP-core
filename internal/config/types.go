package config

type MailerProvider string
type DBProvider string
type FSProvider string
type LLMProvider string
type MORProvider string

// There may be several cases where you use one or the other provider where you use both where it is optional
const (
	MailGmail  MailerProvider = "gmail"
	MailResend MailerProvider = "resend"
)

const (
	DBMongo DBProvider = "mongo"
)

const (
	FSMinio    FSProvider = "minio"
	FSOptional FSProvider = "optional"
)

const (
	LLMOpenAI   LLMProvider = "openai"
	LLMOptional LLMProvider = "optional"
)
