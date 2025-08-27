package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	AppLocation string
	DBUser      string
	DBPass      string
	DBHost      string
	DBPort      string
	DBName      string
	JWTKey      string
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	// Load .env once
	once.Do(func() {
		err := godotenv.Load()

		if err != nil {
			log.Println(".env file not found, using system env")
		}

		cfg = &Config{
			AppPort:     getEnv("APP_PORT", "8080"),
			AppLocation: getEnv("APP_LOCATION", "Asia/Jakarta"),
			DBUser:      getEnv("DB_USER", "root"),
			DBPass:      getEnv("DB_PASS", ""),
			DBHost:      getEnv("DB_HOST", "127.0.0.1"),
			DBPort:      getEnv("DB_PORT", "3306"),
			DBName:      getEnv("DB_NAME", "to_do_list_golang"),
			JWTKey:      getEnv("JWT_KEY", "example_jwt_key"),
		}
	})

	return cfg
}

// Get env with default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
