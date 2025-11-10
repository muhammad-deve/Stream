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

		// Since main.go is inside app/cmd, .env is in app/
		envFilePath := filepath.Join("..", ".env")
		absPath, err := filepath.Abs(envFilePath)
		if err != nil {
			fmt.Println("Env file path error:", err)
		}
		fmt.Printf("Looking for .env file at: %s\n", absPath)

		if err := cleanenv.ReadConfig(envFilePath, instance); err != nil {
			fmt.Printf("❌ Failed to read config from %s: %v\n", envFilePath, err)
			helpText := "Saidoff - Reading project!"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			fmt.Println("⚠️ Application is starting with default config")
		} else {
			fmt.Printf("✅ Successfully loaded config from %s\n", envFilePath)
		}

		if instance.FeaturedChannels != "" {
			fmt.Println("Loaded FeaturedChannels: ✓")
		}
	})
	return instance
}
