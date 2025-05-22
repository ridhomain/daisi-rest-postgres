package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	PgDsn     string
	SecretKey string
}

// LoadConfig initializes Viper, reads defaults, then .env (if not prod), then real env vars.
func LoadConfig() *Config {
	// tell Viper to read environment variables
	viper.AutomaticEnv()
	// allow nested keys like X_API_KEY to map to XApiKey in our struct
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// defaults
	viper.SetDefault("PORT", "3000")
	viper.SetDefault("PG_DSN", "")
	viper.SetDefault("SECRET_KEY", "")

	// load .env in dev (APP_ENV != "production")
	if !isProduction() {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatalf("Error reading .env file: %v", err)
			}
			// missing .env is fine in dev
		}
	}

	return &Config{
		Port:      viper.GetString("PORT"),
		PgDsn:     viper.GetString("PG_DSN"),
		SecretKey: viper.GetString("SECRET_KEY"),
	}
}

func isProduction() bool {
	return strings.EqualFold(os.Getenv("GO_ENV"), "production")
}
