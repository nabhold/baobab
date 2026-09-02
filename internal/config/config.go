package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	HTTPAddress          string
	DatabaseURL          string
	AdminOIDCIssuer      string
	AdminOIDCAudience    string
	WorkloadOIDCIssuer   string
	WorkloadOIDCAudience string
}

func Load() (Config, error) {
	c := Config{
		HTTPAddress:          env("HTTP_ADDRESS", ":8080"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		AdminOIDCIssuer:      os.Getenv("ADMIN_OIDC_ISSUER"),
		AdminOIDCAudience:    env("ADMIN_OIDC_AUDIENCE", "baobab-control-plane"),
		WorkloadOIDCIssuer:   os.Getenv("WORKLOAD_OIDC_ISSUER"),
		WorkloadOIDCAudience: env("WORKLOAD_OIDC_AUDIENCE", "baobab-control-plane"),
	}
	if c.DatabaseURL == "" || c.AdminOIDCIssuer == "" || c.AdminOIDCAudience == "" || c.WorkloadOIDCIssuer == "" || c.WorkloadOIDCAudience == "" {
		return Config{}, errors.New("DATABASE_URL, ADMIN_OIDC_ISSUER, ADMIN_OIDC_AUDIENCE, WORKLOAD_OIDC_ISSUER and WORKLOAD_OIDC_AUDIENCE are required")
	}
	for name, rawIssuer := range map[string]string{"ADMIN_OIDC_ISSUER": c.AdminOIDCIssuer, "WORKLOAD_OIDC_ISSUER": c.WorkloadOIDCIssuer} {
		issuer, err := url.Parse(rawIssuer)
		if err != nil || issuer.Host == "" || (issuer.Scheme != "https" && !localIssuer(issuer)) {
			return Config{}, fmt.Errorf("%s must use HTTPS (HTTP is allowed only for localhost development)", name)
		}
	}
	return c, nil
}

func localIssuer(issuer *url.URL) bool {
	host := strings.ToLower(issuer.Hostname())
	return issuer.Scheme == "http" && (host == "localhost" || host == "127.0.0.1" || host == "::1")
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
