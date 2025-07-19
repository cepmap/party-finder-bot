package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken         string
	DatabaseURL      string
	Port             int
	DBUser           string
	DBPassword       string
	DBName           string
	AdminTelegramIDs []int64
}

func Load(fileName string) *Config {
	// Загружаем переменные из .env файла
	err := godotenv.Load(fileName)
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	config := &Config{
		BotToken:         getEnv("BOT_TOKEN", ""),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		Port:             getEnvAsInt("PORT", 5432),
		DBUser:           getEnv("POSTGRES_USER", ""),
		DBPassword:       getEnv("POSTGRES_PASSWORD", ""),
		DBName:           getEnv("POSTGRES_DB_NAME", "party_finder_bot"),
		AdminTelegramIDs: getEnvAsIntSlice("ADMIN_TELEGRAM_IDS", []int64{}),
	}

	// Проверяем обязательные переменные
	if config.BotToken == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	return config
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как int или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsSlice получает переменную окружения как slice строк (разделенных запятыми)
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// getEnvAsIntSlice получает переменную окружения как slice int64 (разделенных запятыми)
func getEnvAsIntSlice(key string, defaultValue []int64) []int64 {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]int64, 0, len(parts))
		for _, part := range parts {
			if intValue, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64); err == nil {
				result = append(result, intValue)
			}
		}
		return result
	}
	return defaultValue
}
