package config

import (
	"net/smtp"

	gmailsmtp "github.com/GoEnterpricePlatform/goEP-core/internal/gmail-smtp"
	minioClient "github.com/GoEnterpricePlatform/goEP-core/internal/minio"
	mongoClient "github.com/GoEnterpricePlatform/goEP-core/internal/mongo"
	openai "github.com/GoEnterpricePlatform/goEP-core/internal/open-ai"
	resendClient "github.com/GoEnterpricePlatform/goEP-core/internal/resend"

	"github.com/resend/resend-go/v2"
)

type AppClients struct {
	// Mail
	ResendCli *resend.Client
	GmailSmtp smtp.Auth

	// DB
	MongoConn *mongoClient.Data

	// File Storage
	MinioCli *minioClient.MinioClient

	// Optionals services
	OpenaiCli *openai.OpenaiClient
}

func NewClients() *AppClients {
	return &AppClients{}
}

func (ac *AppClients) GetClients(appStack *AppStack, appEnvs *AppEnvs) ( error) {
	switch appStack.Mail {
	case MailResend:
		ac.ResendCli = resendClient.NewResendClient(appEnvs.ResendApiKey)
	case MailGmail:
		ac.GmailSmtp = gmailsmtp.NewGmailSmtpClient(appEnvs.GmailUsername, appEnvs.GmailPass, appEnvs.GmailHost)
	}

	switch appStack.DB {
	case DBMongo:
		ac.MongoConn = mongoClient.New(appEnvs.MongoDBUri)
	}

	switch appStack.FileStg {
	case FSMinio:
		minioCli, err := minioClient.NewClient(appEnvs.MinioEndpoint, appEnvs.MinioAccessKey, appEnvs.MinioSecretKey, appEnvs.MinioUseSSL)
		if err != nil {
			return err
		}
		ac.MinioCli = minioCli
	}

	// Optionals services
	if appEnvs.OpenAiApiKey != "" {
		ac.OpenaiCli = openai.NewOpenAIClient(appEnvs.OpenAiApiKey)
	} 

	return nil
}
