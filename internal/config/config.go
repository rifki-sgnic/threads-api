package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpiryHours int
	EnableDocs     bool
	CORSOrigins    []string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	expiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil || expiry <= 0 {
		return nil, fmt.Errorf("JWT_EXPIRY_HOURS must be a positive integer")
	}

	enableDocs, err := strconv.ParseBool(getEnv("ENABLE_DOCS", "true"))
	if err != nil {
		return nil, errors.New("ENABLE_DOCS must be true or false")
	}

	cfg := &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", ""),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		JWTExpiryHours: expiry,
		EnableDocs:     enableDocs,
		CORSOrigins:    strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
	}

	if len(cfg.JWTSecret) < 32 {
		return nil, errors.New("JWT_SECRET must be set and at least 32 characters")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
