package config

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	FeaturedChannels string `env:"FEATURED_CHANNES"`
	RedisHost        string `env:"REDIS_HOST" env-default:"localhost"`
	RedisPort        string `env:"REDIS_PORT" env-default:"6379"`
	RedisPassword    string `env:"REDIS_PASSWORD" env-default:""`
	RedisDB          int    `env:"REDIS_DB" env-default:"0"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		log.Println("config initializing...")

		instance = &Config{}

		// Try to read from .env file (for local development)
		envFilePath := filepath.Join("..", ".env")
		absPath, err := filepath.Abs(envFilePath)
		if err != nil {
			fmt.Println("Env file path error:", err)
		}
		fmt.Printf("Looking for .env file at: %s\n", absPath)

		// Try to read from file, but don't fail if it doesn't exist
		if err := cleanenv.ReadConfig(envFilePath, instance); err != nil {
			fmt.Printf("⚠️  .env file not found, reading from environment variables\n")
		} else {
			fmt.Printf("✅ Successfully loaded config from %s\n", envFilePath)
		}

		// Always read from environment variables (overrides file values)
		if err := cleanenv.ReadEnv(instance); err != nil {
			fmt.Printf("❌ Failed to read environment variables: %v\n", err)
		} else {
			fmt.Println("✅ Environment variables loaded")
		}

		// Display loaded configuration
		fmt.Printf("Configuration:\n")
		fmt.Printf("  REDIS_HOST: %s\n", instance.RedisHost)
		fmt.Printf("  REDIS_PORT: %s\n", instance.RedisPort)
		fmt.Printf("  REDIS_DB: %d\n", instance.RedisDB)
		fmt.Printf("  REDIS_PASSWORD: %s\n", maskPassword(instance.RedisPassword))
		if instance.FeaturedChannels != "" {
			fmt.Println("  FEATURED_CHANNES: ✓ =======================================================>", instance.FeaturedChannels)
		}
	})
	return instance
}

func maskPassword(password string) string {
	if password == "" {
		return "(empty)"
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}
