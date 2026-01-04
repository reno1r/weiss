package config

import "github.com/spf13/viper"

type Config struct {
	AppName  string `mapstructure:"APP_NAME"`
	AppHost  string `mapstructure:"APP_HOST"`
	AppPort  string `mapstructure:"APP_PORT"`
	AppDebug bool   `mapstructure:"APP_DEBUG"`

	DatabaseHost                string `mapstructure:"DB_HOST"`
	DatabasePort                int    `mapstructure:"DB_PORT"`
	DatabaseName                string `mapstructure:"DB_NAME"`
	DatabaseUser                string `mapstructure:"DB_USER"`
	DatabasePassword            string `mapstructure:"DB_PASSWORD"`
	DatabaseSSL                 string `mapstructure:"DB_SSL"`
	DatabaseMaxConnections      int    `mapstructure:"DB_MAX_CONNECTIONS"`
	DatabaseMaxIdleConnections  int    `mapstructure:"DB_MAX_IDLE_CONNECTIONS"`
	DatabaseConnectionTimeoutMs int    `mapstructure:"DB_CONNECTION_TIMEOUT_MS"`

	BcryptCost         int    `mapstructure:"BCRYPT_COST"`
	CorsAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`

	JwtSecret         string `mapstructure:"JWT_SECRET"`
	JwtIssuer         string `mapstructure:"JWT_ISSUER"`
	JwtAccessExpires  string `mapstructure:"JWT_ACCESS_EXPIRES_IN"`
	JwtRefreshExpires string `mapstructure:"JWT_REFRESH_EXPIRES_IN"`
}

var config *Config

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
