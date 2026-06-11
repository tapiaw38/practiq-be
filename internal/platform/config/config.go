package config

type (
	ServerConfig struct {
		AppName     string
		Port        string
		GinMode     string
		JWTSecret   string
		FrontendURL string
		AuthAPIURL  string
	}

	DatabaseConfig struct {
		DatabaseURL string
	}

	S3Config struct {
		AWSRegion          string
		AWSAccessKeyID     string
		AWSSecretAccessKey string
		AWSSessionToken    string
		AWSBucket          string
		AWSEndpoint        string
	}

	Config struct {
		ServerConfig   ServerConfig
		DatabaseConfig DatabaseConfig
		S3Config       S3Config
	}
)
