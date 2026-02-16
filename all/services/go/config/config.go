package config

import (
	"os"
)

type Config struct {
	Env        string
	HTTPAddr   string
	MasterKey  string
}

func Load() *Config {
	return &Config{
		Env:       getEnv("APP_ENV", "dev"),
		HTTPAddr:  getEnv("HTTP_ADDR", "127.0.0.1:8081"),
		MasterKey: getEnv("APP_MASTER_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
