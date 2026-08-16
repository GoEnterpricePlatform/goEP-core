package config

import (
	"cmp"
	"log"
	"os"
	"strings"
	"time"
)

type AppEnvs struct {
	// MongoDB
	MongoDBUri  string
	MongoInitDB string

	// Minio
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioUseSSL     bool
	MinioBucketName string

	// Resend
	ResendApiKey    string
	ResendEmailFrom string

	// GmailSmtp
	GmailUsername string
	GmailPass     string
	GmailFrom     string
	GmailHost     string
	GmailAddr     string

	// OpenAI
	OpenAiApiKey string

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

func (ae *AppEnvs) Load(appStack *AppStack) {
	// Mongo DB
	mongoInitDB := cmp.Or(os.Getenv("MONGO_INITDB_DATABASE"), "goep-core")

	// Minio
	minioBucketName := cmp.Or(os.Getenv("MINIO_BUCKET_NAME"), "goep-core")
	useSSL := mustGetEnv("MINIO_SECURE")

	var useSSLbool bool
	if useSSL == "true" || useSSL == "yes" {
		useSSLbool = true
	} else {
		useSSLbool = false
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

	openaiApiKey := cmp.Or(os.Getenv("OPENAI_API_KEY"), "")

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
	// Resend
	var resendApiKey string
	var resendEmailFrom string

	// GmailSmtp
	var gmailUsername string
	var gmailPass string
	var gmailFrom string
	var gmailHost string
	var gmailAddr string

	switch appStack.Mail {
	case MailResend:
		resendApiKey = mustGetEnv("RESEND_API_KEY")
		resendEmailFrom = mustGetEnv("EMAIL_FROM")
	case MailGmail:
		gmailUsername = mustGetEnv("GMAIL_USERNAME")
		gmailPass = mustGetEnv("GMAIL_PASS")
		gmailFrom = cmp.Or(os.Getenv("GMAIL_FROM"), gmailUsername)
		gmailHost = cmp.Or(os.Getenv("GMAIL_HOST"), "smtp.gmail.com")
		gmailAddr = cmp.Or(os.Getenv("GMAIL_ADDR"), "smtp.gmail.com:587")
	}

	ae.MongoDBUri = mustGetEnv("MONGO_DB_URI")
	ae.MongoInitDB = mongoInitDB
	ae.MinioEndpoint = mustGetEnv("MINIO_ENDPOINT")
	ae.MinioAccessKey = mustGetEnv("MINIO_ACCESS_KEY")
	ae.MinioSecretKey = mustGetEnv("MINIO_SECRET_KEY")
	ae.MinioUseSSL = useSSLbool
	ae.MinioBucketName = minioBucketName
	ae.ResendApiKey = resendApiKey
	ae.ResendEmailFrom = resendEmailFrom
	ae.AppName = appName
	ae.AppDataPath = appDataPath
	ae.GmailUsername = gmailUsername
	ae.GmailPass = gmailPass
	ae.GmailFrom = gmailFrom
	ae.GmailHost = gmailHost
	ae.GmailAddr = gmailAddr
	ae.OpenAiApiKey = openaiApiKey
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
