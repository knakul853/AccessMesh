package config

import "os"

type Config struct {
	MongoURI    string
	JWTSecret   string
	Environment string
}

func Load() *Config {
	return &Config{
		MongoURI:    getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		Environment: getEnvOrDefault("ENV", "development"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
