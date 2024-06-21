package config

import "github.com/spf13/viper"

// Stores the configuration for the application.
// The values are read by viper from the config file or environment variables.
type Config struct {
	ENV            string `mapstructure:"ENVIRONMENT"`
	Port           string `mapstructure:"HTTP_PORT"`
	DBURL          string `mapstructure:"DATABASE_URL"`
	RodURL         string `mapstructure:"ROD_BROWSER_URL"`
	AdminSub       string `mapstructure:"ADMIN_SUB"`
	SentryDSN      string `mapstructure:"SENTRY_DSN"`
	RedisURL       string `mapstructure:"REDIS_URL"`
	Version        string `mapstructure:"VERSION"`
	SearchURL      string `mapstructure:"OPENSEARCH_URL"`
	ClerkSecretKey string `mapstructure:"CLERK_SECRET_KEY"`
	KafkaURL       string `mapstructure:"KAFKA_URL"`
	KafkaUsername  string `mapstructure:"KAFKA_USERNAME"`
	KafkaPassword  string `mapstructure:"KAFKA_PASSWORD"`
}

// Reads the configuration from the config file or environment variables.
// Returns the configuration as Config struct or an error.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
