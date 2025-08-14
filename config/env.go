package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	User                   string
	Password               string
	Host                   string
	Port                   int
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	portStr := getEnv("DB_PORT", "")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 1433
	}

	return Config{
		User:                   getEnv("DB_USER", ""),
		Password:               getEnv("DB_PASSWORD", ""),
		Host:                   getEnv("DB_HOST", ""),
		Port:                   port,
		DBName:                 getEnv("DB_NAME", ""),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return result
	}
	return fallback
}
