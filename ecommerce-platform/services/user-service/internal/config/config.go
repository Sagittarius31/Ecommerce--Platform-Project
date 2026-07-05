package config

import "os"

type Config struct {
	Port        string
	Env         string
	Version     string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		Version:     getEnv("VERSION", "1.0.0"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user_svc:secret123@localhost:5432/users_db?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "change-this-in-production"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
