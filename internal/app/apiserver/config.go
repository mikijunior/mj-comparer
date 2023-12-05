package apiserver

import (
	"fmt"
	"os"
)

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
	SessionKey  string
}

func NewConfig() *Config {
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		BindAddr: fmt.Sprintf(":%s", os.Getenv("PORT")),
		LogLevel: os.Getenv("LOG_LEVEL"),
		SessionKey: os.Getenv("SESSION_KEY"),
	}
}