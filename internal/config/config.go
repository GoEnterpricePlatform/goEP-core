package config

import (
	"cmp"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
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
	ResendApiKey string
	EmailFrom    string
	AppName      string

	// Jwt
	JWTAccessSecret           string
	JWTRefreshSecret          string
	JWTIssuer                 string
	JWTAccessExpIn            time.Duration
	JWTRefreshExpIn           time.Duration
	JWTRefreshRememberMeExpIn time.Duration

	// App
	Port           string
	AppEnv         string
	AllowedOrigins []string

	// Templates
	ApiBaseUrl string
}

func Load() *Config {
	// Mongo DB
	mongoInitDB := cmp.Or(os.Getenv("MONGO_INITDB_DATABASE"), "auth-tmpl")

	// Minio
	minioBucketName := cmp.Or(os.Getenv("MINIO_BUCKET_NAME"), "auth-tmpl")
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

	// App
	port := cmp.Or(os.Getenv("HTTP_SERVER_PORT"), "8000")
	appEnv := cmp.Or(os.Getenv("APP_ENV"), "dev")
	apiBaseUrl := cmp.Or(os.Getenv("API_BASE_URL"), "http://localhost:"+port)

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	return &Config{
		MongoDBUri:                mustGetEnv("MONGO_DB_URI"),
		MongoInitDB:               mongoInitDB,
		MinioEndpoint:             mustGetEnv("MINIO_ENDPOINT"),
		MinioAccessKey:            mustGetEnv("MINIO_ACCESS_KEY"),
		MinioSecretKey:            mustGetEnv("MINIO_SECRET_KEY"),
		MinioUseSSL:               useSSLbool,
		MinioBucketName:           minioBucketName,
		ResendApiKey:              mustGetEnv("RESEND_API_KEY"),
		EmailFrom:                 mustGetEnv("EMAIL_FROM"),
		AppName:                   mustGetEnv("APP_NAME"),
		JWTAccessSecret:           mustGetEnv("JWT_ACCESS_TOKEN"),
		JWTRefreshSecret:          mustGetEnv("JWT_REFRESH_TOKEN"),
		JWTIssuer:                 mustGetEnv("JWT_ISS"),
		JWTAccessExpIn:            accessDur,
		JWTRefreshExpIn:           refreshDur,
		JWTRefreshRememberMeExpIn: refreshRememberMeDur,
		Port:                      port,
		AppEnv:                    appEnv,
		ApiBaseUrl:                apiBaseUrl,
		AllowedOrigins:            allowedOrigins,
	}
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return val
}
