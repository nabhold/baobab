package config

import (
	"errors"
	"os"
)

type Config struct{ HTTPAddress, DatabaseURL, AdminToken string }

func Load() (Config, error) {
	c := Config{env("HTTP_ADDRESS", ":8080"), os.Getenv("DATABASE_URL"), os.Getenv("ADMIN_API_TOKEN")}
	if c.DatabaseURL == "" || len(c.AdminToken) < 32 {
		return Config{}, errors.New("DATABASE_URL and ADMIN_API_TOKEN (minimum 32 characters) are required")
	}
	return c, nil
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
