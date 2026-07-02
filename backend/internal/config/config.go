package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                  string
	DatabaseURL           string
	FrontendURL           string
	BackendURL            string
	AuthentikIssuerURL    string
	AuthentikClientID     string
	AuthentikClientSecret string
	AuthentikRedirectURL  string
	CORSOrigins           []string
	DevAuth               bool
	AutoImportExcel       bool
	ImportExcelPath       string
}

func Load() Config {
	frontend := env("FRONTEND_URL", "http://localhost:5173")
	return Config{
		Port:                  env("PORT", "8080"),
		DatabaseURL:           env("DATABASE_URL", "postgres://airclean:airclean@localhost:5432/airclean?sslmode=disable"),
		FrontendURL:           frontend,
		BackendURL:            env("BACKEND_URL", "http://localhost:8080"),
		AuthentikIssuerURL:    os.Getenv("AUTHENTIK_ISSUER_URL"),
		AuthentikClientID:     os.Getenv("AUTHENTIK_CLIENT_ID"),
		AuthentikClientSecret: os.Getenv("AUTHENTIK_CLIENT_SECRET"),
		AuthentikRedirectURL:  os.Getenv("AUTHENTIK_REDIRECT_URL"),
		CORSOrigins:           split(env("CORS_ORIGINS", frontend)),
		DevAuth:               boolEnv("DEV_AUTH", true),
		AutoImportExcel:       boolEnv("AUTO_IMPORT_EXCEL", false),
		ImportExcelPath:       os.Getenv("IMPORT_EXCEL_PATH"),
	}
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func boolEnv(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return v
}

func split(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if v := strings.TrimSpace(part); v != "" {
			out = append(out, v)
		}
	}
	return out
}
