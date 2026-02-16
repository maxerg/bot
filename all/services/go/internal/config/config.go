package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
        Env                string
        HTTPAddr           string
        CORSAllowedOrigins []string
        MasterKey          string
        TelegramBotToken   string
        DBDsn              string
}



func MustLoad() Config {
	cfg := Config{
		Env:                getenv("APP_ENV", "dev"),
		HTTPAddr:           getenv("HTTP_ADDR", ":8081"),
		MasterKey:          getenv("APP_MASTER_KEY", ""),
		TelegramBotToken:   getenv("TELEGRAM_BOT_TOKEN", ""),
		DBDsn:              getenv("DB_DSN", ""),
		CORSAllowedOrigins: splitCSV(getenv("CORS_ALLOWED_ORIGINS", "")),

	}
	if len(cfg.MasterKey) < 32 {
		log.Println("warning: APP_MASTER_KEY should be at least 32 bytes")
	}
	return cfg
}

func getenv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

func splitCSV(v string) []string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
