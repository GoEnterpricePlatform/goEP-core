package config

import (
	"cmp"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func (ae *AppEnvs) Load() {
	// Database
	dbProviderValue := cmp.Or(os.Getenv("DB_PROVIDER"), "mongo")
	var dbProvider DBProvider
	dbProvider = DBProvider(dbProviderValue)

	// Mongo
	var mongoInitDB string
	var mongoDBUri string

	switch dbProvider {
	case DBMongo:
		mongoInitDB = cmp.Or(os.Getenv("MONGO_INITDB_DATABASE"), "goep-core")
		mongoDBUri = mustGetEnv("MONGO_DB_URI")

	default:
		log.Fatalf("invalid DB_PROVIDER: %q", dbProviderValue)
	}

	// File storage
	fileStorageProviderValue := cmp.Or(os.Getenv("FILE_STORAGE_PROVIDER"), "optional")
	var fileStorageProvider FSProvider
	fileStorageProvider = FSProvider(fileStorageProviderValue)

	// Minio
	var useSSLbool bool
	var minioBucketName string
	var minioEndpoint string
	var minioAccessKey string
	var minioSecretKey string

	switch fileStorageProvider {

	case FSMinio:
		minioBucketName = cmp.Or(os.Getenv("MINIO_BUCKET_NAME"), "goep-core")
		useSSL := mustGetEnv("MINIO_SECURE")

		if useSSL == "true" || useSSL == "yes" {
			useSSLbool = true
		} else {
			useSSLbool = false
		}

		minioEndpoint = mustGetEnv("MINIO_ENDPOINT")
		minioAccessKey = mustGetEnv("MINIO_ACCESS_KEY")
		minioSecretKey = mustGetEnv("MINIO_SECRET_KEY")

	case FSOptional:

	default:
		log.Fatalf("invalid FILE_STORAGE_PROVIDER: %q", fileStorageProviderValue)
	}

	// We don't know if one or the other MoR or both will be used, so I will validate it separately and not switch
	isEnableMorPaddle, err := strconv.ParseBool(cmp.Or(os.Getenv("IS_ENABLE_MOR_PADDLE"), "false"))
	if err != nil {
		log.Fatal("IS_ENABLE_MOR_PADDLE must be a boolean")
	}

	// paddle
	var paddleApiKey string
	paddleEnvironment := cmp.Or(os.Getenv("PADDLE_ENVIRONMENT"), "sandbox")
	var paddleWebhookSecret string

	if isEnableMorPaddle {
		paddleApiKey = mustGetEnv("PADDLE_API_KEY")
		paddleWebhookSecret = mustGetEnv("PADDLE_WEBHOOK_SECRET_KEY")
		if paddleEnvironment != "sandbox" && paddleEnvironment != "production" {
			log.Fatalf("PADDLE_ENVIRONMENT must be either sandbox or production")
		}
	}

	// Auth - tokens
	accessExp := cmp.Or(os.Getenv("JWT_ACCESS_EXP_IN"), "15m")
	refreshExp := cmp.Or(os.Getenv("JWT_REFRESH_EXP_IN"), "168h")
	refreshExpRememberMe := cmp.Or(os.Getenv("JWT_REFRESH_EXP_IN_REMEMBER"), "720h")

	accessDur, err := time.ParseDuration(accessExp)
	if err != nil {
		log.Fatalf("Invalid JWT_ACCESS_EXP_IN format: %v", err)
	}

	refreshDur, err := time.ParseDuration(refreshExp)
	if err != nil {
		log.Fatalf("Invalid JWT_REFRESH_EXP_IN format: %v", err)
	}

	refreshRememberMeDur, err := time.ParseDuration(refreshExpRememberMe)
	if err != nil {
		log.Fatalf("Invalid JWT_REFRESH_EXP_IN_REMEMBER format: %v", err)
	}

	// LLM provider
	llmProviderValue := cmp.Or(os.Getenv("LLM_PROVIDER"), "optional")
	var llmProvider LLMProvider
	llmProvider = LLMProvider(llmProviderValue)

	// OpenAI
	var openaiApiKey string

	switch llmProvider {
	case LLMOpenAI:
		openaiApiKey = mustGetEnv("OPENAI_API_KEY")

	case LLMOptional:

	default:
		log.Fatalf("invalid LLM_PROVIDER: %q", llmProviderValue)
	}

	// Cookies
	// JWT_ACCESS_COOKIE_EXP_IN defines how long the access token cookie
	// stays in the browser. It must be greater than JWT_ACCESS_EXP_IN
	// so the backend can detect an expired JWT and refresh it before
	// the browser removes the cookie automatically.
	jwtAccessCookieExpIn := cmp.Or(os.Getenv("JWT_ACCESS_COOKIE_EXP_IN"), "20m")
	jwtAccessCookieDur, err := time.ParseDuration(jwtAccessCookieExpIn)
	if err != nil {
		log.Fatalf("Invalid JWT_ACCESS_EXP_IN format: %v", err)
	}

	// App
	port := cmp.Or(os.Getenv("HTTP_SERVER_PORT"), "8000")
	appEnv := cmp.Or(os.Getenv("APP_ENV"), "dev")
	apiBaseUrl := cmp.Or(os.Getenv("API_BASE_URL"), "http://localhost:"+port)

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	// "appname" appears in the email message
	appName := mustGetEnv("APP_NAME")

	appDataPath := appName + "_data"

	// mailer
	mailerProvider := MailerProvider(cmp.Or(os.Getenv("MAILER_PROVIDER"), "gmail"))

	// Resend
	var resendApiKey string
	var resendEmailFrom string

	// GmailSmtp
	var gmailUsername string
	var gmailPass string
	var gmailFrom string
	var gmailHost string
	var gmailAddr string

	switch mailerProvider {
	case MailGmail:
		gmailUsername = mustGetEnv("GMAIL_USERNAME")
		gmailPass = mustGetEnv("GMAIL_PASS")
		gmailFrom = cmp.Or(os.Getenv("GMAIL_FROM"), gmailUsername)
		gmailHost = cmp.Or(os.Getenv("GMAIL_HOST"), "smtp.gmail.com")
		gmailAddr = cmp.Or(os.Getenv("GMAIL_ADDR"), "smtp.gmail.com:587")
	case MailResend:
		resendApiKey = mustGetEnv("RESEND_API_KEY")
		resendEmailFrom = mustGetEnv("EMAIL_FROM")
	default:
		log.Fatalf("invalid MAILER_PROVIDER: %q", mailerProvider)
	}

	ae.DBProvider = dbProvider
	ae.MongoDBUri = mongoDBUri
	ae.MongoInitDB = mongoInitDB
	ae.FileStorageProvider = fileStorageProvider
	ae.MinioEndpoint = minioEndpoint
	ae.MinioAccessKey = minioAccessKey
	ae.MinioSecretKey = minioSecretKey
	ae.MinioUseSSL = useSSLbool
	ae.MinioBucketName = minioBucketName
	ae.MailerProvider = mailerProvider
	ae.ResendApiKey = resendApiKey
	ae.ResendEmailFrom = resendEmailFrom
	ae.AppName = appName
	ae.AppDataPath = appDataPath
	ae.GmailUsername = gmailUsername
	ae.GmailPass = gmailPass
	ae.GmailFrom = gmailFrom
	ae.GmailHost = gmailHost
	ae.GmailAddr = gmailAddr
	ae.LLMProvider = llmProvider
	ae.OpenAiApiKey = openaiApiKey
	ae.IsEnableMorPaddle = isEnableMorPaddle
	ae.PaddleApiKey = paddleApiKey
	ae.PaddleEnvironment = paddleEnvironment
	ae.PaddleWebhookSecret = paddleWebhookSecret
	ae.JWTAccessSecret = mustGetEnv("JWT_ACCESS_TOKEN")
	ae.JWTRefreshSecret = mustGetEnv("JWT_REFRESH_TOKEN")
	ae.JWTIssuer = mustGetEnv("JWT_ISS")
	ae.JWTAccessExpIn = accessDur
	ae.JWTRefreshExpIn = refreshDur
	ae.JWTRefreshRememberMeExpIn = refreshRememberMeDur
	ae.JwtAccessCookieDur = jwtAccessCookieDur
	ae.Port = port
	ae.AppEnv = appEnv
	ae.ApiBaseUrl = apiBaseUrl
	ae.AllowedOrigins = allowedOrigins
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return val
}
