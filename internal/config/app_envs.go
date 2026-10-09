package config

import (
	"time"
)

type AppEnvs struct {
	DBProvider DBProvider

	// MongoDB
	MongoDBUri  string
	MongoInitDB string

	// File Storage provider
	FileStorageProvider FSProvider

	// Minio
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioUseSSL     bool
	MinioBucketName string

	// Mail provider
	MailerProvider MailerProvider

	// Resend
	ResendApiKey    string
	ResendEmailFrom string

	// GmailSmtp
	GmailUsername string
	GmailPass     string
	GmailFrom     string
	GmailHost     string
	GmailAddr     string

	// LLM provider
	LLMProvider LLMProvider

	// OpenAI
	OpenAiApiKey string

	// payments MoR
	IsEnableMorPaddle   bool
	PaddleApiKey        string
	PaddleEnvironment   string
	PaddleWebhookSecret string

	// Jwt
	JWTAccessSecret           string
	JWTRefreshSecret          string
	JWTIssuer                 string
	JWTAccessExpIn            time.Duration
	JWTRefreshExpIn           time.Duration
	JWTRefreshRememberMeExpIn time.Duration
	JwtAccessCookieDur        time.Duration

	// App
	Port           string
	AppEnv         string
	AllowedOrigins []string
	AppName        string // send email resend, gmail
	AppDataPath    string

	// Templates
	ApiBaseUrl string
}

func NewAppEnvs() *AppEnvs {
	return &AppEnvs{}
}
