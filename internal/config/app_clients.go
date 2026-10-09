package config

import (
	"net/smtp"

	gmailsmtp "github.com/GoEnterpricePlatform/goEP-core/internal/gmail-smtp"
	minioClient "github.com/GoEnterpricePlatform/goEP-core/internal/minio"
	mongoClient "github.com/GoEnterpricePlatform/goEP-core/internal/mongo"
	openai "github.com/GoEnterpricePlatform/goEP-core/internal/open-ai"
	paddleClient "github.com/GoEnterpricePlatform/goEP-core/internal/paddle"
	resendClient "github.com/GoEnterpricePlatform/goEP-core/internal/resend"

	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppClients struct {
	// Mail
	ResendCli *resend.Client
	GmailSmtp smtp.Auth

	// DB
	MongoConn *mongoClient.Data
	DB        *mongo.Database

	// File Storage
	MinioCli *minioClient.MinioClient

	// Optionals services
	OpenaiCli *openai.OpenaiClient

	// optional paddle
	PaddleCli *paddleClient.PaddleClient
}

func NewClients() *AppClients {
	return &AppClients{}
}

func (ac *AppClients) GetClients(appEnvs *AppEnvs) error {
	switch appEnvs.MailerProvider {
	case MailResend:
		ac.ResendCli = resendClient.NewResendClient(appEnvs.ResendApiKey)
	case MailGmail:
		ac.GmailSmtp = gmailsmtp.NewGmailSmtpClient(appEnvs.GmailUsername, appEnvs.GmailPass, appEnvs.GmailHost)
	}

	switch appEnvs.DBProvider {
	case DBMongo:
		ac.MongoConn = mongoClient.New(appEnvs.MongoDBUri)
	}

	switch appEnvs.FileStorageProvider {
	case FSMinio:
		minioCli, err := minioClient.NewClient(appEnvs.MinioEndpoint, appEnvs.MinioAccessKey, appEnvs.MinioSecretKey, appEnvs.MinioUseSSL)
		if err != nil {
			return err
		}
		ac.MinioCli = minioCli
	}

	switch appEnvs.LLMProvider {
	case LLMOpenAI:
		ac.OpenaiCli = openai.NewOpenAIClient(appEnvs.OpenAiApiKey)
	}

	if appEnvs.IsEnableMorPaddle {
		paddleCli, err := paddleClient.NewPaddleClient(appEnvs.PaddleApiKey, appEnvs.PaddleEnvironment)
		if err != nil {
			return err
		}
		ac.PaddleCli = paddleCli
	}
	return nil
}

// Initialize the infrastructure in the services that we are using and that are strictly
// necessary for all modules; if it is only for a specific module, it should go in the
// initializer folder within your module.
func (ac *AppClients) InitializeServices(appEnvs *AppEnvs) error {

	switch appEnvs.DBProvider {
	case DBMongo:
		// - MongoDB
		ac.DB = ac.MongoConn.DB.Database(appEnvs.MongoInitDB)
		ac.MongoConn.Ping()
	}

	switch appEnvs.FileStorageProvider {
	case FSMinio:
		err := ac.MinioCli.CreateStorage(appEnvs.MinioBucketName)
		if err != nil {
			return err
		}
	}
	return nil
}
