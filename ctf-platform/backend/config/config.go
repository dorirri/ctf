package config

import "os"

type Config struct {
	Port      string
	DBURL     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBURL:     getEnv("DATABASE_URL", ""),
		JWTSecret: getEnv("JWT_SECRET", "changeme"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
