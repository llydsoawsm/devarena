package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                     string
	AppEnv                   string
	ClientURL                string
	DBDriver                 string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	DBSSLMode                string
	SQLitePath               string
	RedisAddr                string
	RedisPassword            string
	JWTSecret                string
	JWTRefreshSecret         string
	JWTExpirationMinutes     int
	JWTRefreshExpirationDays int
	DockerHost               string
	JudgeMaxWorkers          int
	JudgeTimeoutSeconds      int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Port:                     getEnv("PORT", "8080"),
		AppEnv:                   getEnv("APP_ENV", "development"),
		ClientURL:                getEnv("CLIENT_URL", "http://localhost:3000"),
		DBDriver:                 getEnv("DB_DRIVER", "sqlite"),
		DBHost:                   getEnv("DB_HOST", "127.0.0.1"),
		DBPort:                   getEnv("DB_PORT", "5432"),
		DBUser:                   getEnv("DB_USER", "devcode"),
		DBPassword:               getEnv("DB_PASSWORD", "devcode_secret_password"),
		DBName:                   getEnv("DB_NAME", "devcode"),
		DBSSLMode:                getEnv("DB_SSLMODE", "disable"),
		SQLitePath:               getEnv("SQLITE_PATH", "devcode.db"),
		RedisAddr:                getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:            getEnv("REDIS_PASSWORD", ""),
		JWTSecret:                getEnv("JWT_SECRET", "super_secret_devcode_jwt_access_key_change_in_production_32b"),
		JWTRefreshSecret:         getEnv("JWT_REFRESH_SECRET", "super_secret_devcode_jwt_refresh_key_change_in_prod"),
		JWTExpirationMinutes:     getEnvInt("JWT_EXPIRATION_MINUTES", 60),
		JWTRefreshExpirationDays: getEnvInt("JWT_REFRESH_EXPIRATION_DAYS", 7),
		DockerHost:               getEnv("DOCKER_HOST", "npipe:////./pipe/docker_engine"),
		JudgeMaxWorkers:          getEnvInt("JUDGE_MAX_WORKERS", 4),
		JudgeTimeoutSeconds:      getEnvInt("JUDGE_TIMEOUT_SECONDS", 5),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
