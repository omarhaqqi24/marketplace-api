package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	AppPort string
	AppEnv  string
}

func Load() *Config {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load config file: $v", err)
	}

	return &Config{
		AppName: viper.GetString("APP_NAME"),
		AppPort: viper.GetString("APP_PORT"),
		AppEnv:  viper.GetString("APP_ENV"),
	}
}
