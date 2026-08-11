package config

import "time"

type Config struct {
	Debug    bool   `mapstructure:"DEBUG"`
	LogLevel string `mapstructure:"LOG_LEVEL"`

	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`

	DatabaseURL string `mapstructure:"DATABASE_URL"`

	UploadsDirPath string `mapstructure:"UPLOADS_DIR_PATH"`

	// CORSAllowedOrigins is a comma separated list of origins allowed to call
	// the API from a browser. Empty falls back to the local dev origins.
	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`

	SessionExpiredAt time.Duration `mapstructure:"SESSION_EXPIRED_AT"`

	EncodedSessionHashKey  string `mapstructure:"ENCODED_SESSION_HASH_KEY"`
	EncodedSessionBlockKey string `mapstructure:"ENCODED_SESSION_BLOCK_KEY"`

	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

var Configuration = &Config{}
