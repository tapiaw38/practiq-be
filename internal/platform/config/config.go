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

	Config struct {
		ServerConfig   ServerConfig
		DatabaseConfig DatabaseConfig
	}
)
