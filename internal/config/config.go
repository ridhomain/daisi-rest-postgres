package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	PgDsn     string
	SecretKey string
	XApiKey   string
}

// LoadConfig reads environment variables from .env and returns a Config
func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	return &Config{
		Port:      viper.GetString("PORT"),
		PgDsn:     viper.GetString("PG_DSN"),
		SecretKey: viper.GetString("SECRET_KEY"),
	}
}
