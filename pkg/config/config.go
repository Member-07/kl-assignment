package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	MysqlConfig DB
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	cfg := Config{
		Port: env("PORT", ":8080"),
		MysqlConfig: DB{
			DBHost:     env("DB_HOST", "127.0.0.1"),
			DBPort:     env("DB_PORT", "3306"),
			DBUser:     env("DB_USER", "appuser"),
			DBPassword: env("DB_PASSWORD", "my-secret-password"),
			DBName:     env("DB_NAME", "kl-assignment-db"),
		},
	}

	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(v string) (time.Duration, error) {
	return time.ParseDuration(v)
}
