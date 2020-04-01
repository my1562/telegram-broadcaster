package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TGToken     string
	RabbitmqURL string
}

func NewConfig() *Config {

	envFile := ".env"

	injectedEnvFile := os.Getenv("ENV_FILE")
	if injectedEnvFile != "" {
		envFile = injectedEnvFile
	}

	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Println(err)
	}

	config := &Config{
		TGToken:     os.Getenv("TELEGRAM_APITOKEN"),
		RabbitmqURL: os.Getenv("RABBITMQ_URL"),
	}

	return config
}
