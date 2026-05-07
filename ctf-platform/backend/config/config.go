package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv        string
	Port          string
	DBURL         string
	JWTSecret     string
	AdminUsername string
	AdminEmail    string
	AdminPassword string
}

func Load() *Config {
	cfg := &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		Port:          getEnv("PORT", "8080"),
		DBURL:         getFirstEnv([]string{"DATABASE_URL", "DB_URL"}, ""),
		JWTSecret:     getEnv("JWT_SECRET", "changeme"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminEmail:    getEnv("ADMIN_EMAIL", ""),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}

	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) Validate() error {
	if c.DBURL == "" {
		return fmt.Errorf("DATABASE_URL or DB_URL is required")
	}

	if c.AppEnv != "development" && c.JWTSecret == "changeme" {
		return fmt.Errorf("JWT_SECRET must be set outside development")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getFirstEnv(keys []string, fallback string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return fallback
}
