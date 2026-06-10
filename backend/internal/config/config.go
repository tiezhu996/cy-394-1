package config

import "os"

type Config struct {
	DatabaseDSN string
	JWTSecret   string
	ServerPort  string
}

func Load() Config {
	return Config{
		DatabaseDSN: env("DATABASE_DSN", "host=localhost user=fitness password=fitness_pass dbname=fitness port=5432 sslmode=disable TimeZone=Asia/Shanghai"),
		JWTSecret:   env("JWT_SECRET", "change-me"),
		ServerPort:  env("SERVER_PORT", "8080"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
