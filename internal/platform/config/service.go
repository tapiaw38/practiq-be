package config

import "os"

var configService *Config

func InitConfigService() {
	configService = &Config{
		ServerConfig: ServerConfig{
			AppName:     getEnv("APP_NAME", "practiq-be"),
			Port:        getEnv("PORT", "8083"),
			GinMode:     getEnv("GIN_MODE", "debug"),
			JWTSecret:   getEnv("JWT_SECRET", "secret"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5174"),
			AuthAPIURL:  getEnv("AUTH_API_URL", "http://localhost:8082"),
		},
		DatabaseConfig: DatabaseConfig{
			DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:54323/practiq-db?sslmode=disable"),
		},
		S3Config: S3Config{
			AWSRegion:          getEnv("AWS_REGION", ""),
			AWSBucket:          getEnv("AWS_BUCKET", ""),
			AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			AWSSessionToken:    getEnv("AWS_SESSION_TOKEN", ""),
			AWSEndpoint:        getEnv("AWS_ENDPOINT", ""),
		},
	}
}

func GetConfigService() *Config {
	return configService
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
