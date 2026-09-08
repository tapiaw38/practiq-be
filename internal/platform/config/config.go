package config

type (
	ServerConfig struct {
		AppName     string
		Port        string
		GinMode     string
		JWTSecret   string
		FrontendURL string
		AuthAPIURL  string
		// PaymentsURL is reachable only inside the compose network: the
		// payments service is never published, and this is its only client.
		PaymentsURL    string
		PaymentsAPIKey string
		// MercadoPagoPublicKey is served to the browser so it can turn card
		// details into a token without them passing through us. It is public
		// by design; the secret half never leaves the payments service.
		MercadoPagoPublicKey string
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
