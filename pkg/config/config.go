package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	_ = godotenv.Load()
}

func Get(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		n, err := strconv.Atoi(value)
		if err == nil {
			return n
		}
	}
	return defaultValue
}