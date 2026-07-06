package config

import "os"

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
	JWTSecret   string
	RedisURL    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user_svc:secret123@localhost:5432/users_db?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-this-in-production"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
	}
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}
