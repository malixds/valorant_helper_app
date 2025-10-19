package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	WebhookURL       string
	Port             string
	NgrokURL         string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	ValorantAPIKey   string
}

func LoadConfig() *Config {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		WebhookURL:       getEnv("WEBHOOK_URL", ""),
		Port:             getEnv("PORT", "8080"),
		NgrokURL:         getEnv("NGROK_URL", ""),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "valorant_user"),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBName:           getEnv("DB_NAME", "valorant_db"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		ValorantAPIKey:   getEnv("VALORANT_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
