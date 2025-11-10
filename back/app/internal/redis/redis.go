package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gitlab.yurtal.tech/company/blitz/business-card/back/internal/config"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient creates a new Redis client instance
func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// GenerateURLToken generates a UUID token for the given URL and stores it in Redis with 2-hour expiration
func (r *RedisClient) GenerateURLToken(url string) (string, error) {
	// Generate UUID
	token := uuid.New().String()

	// Store URL with token as key, expires in 2 hours
	expiration := 2 * time.Hour
	err := r.client.Set(r.ctx, token, url, expiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store URL in Redis: %w", err)
	}

	return token, nil
}

// GetURLByToken retrieves the URL associated with the given token
func (r *RedisClient) GetURLByToken(token string) (string, error) {
	url, err := r.client.Get(r.ctx, token).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("token not found or expired")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve URL from Redis: %w", err)
	}

	return url, nil
}

// DeleteToken removes a token from Redis
func (r *RedisClient) DeleteToken(token string) error {
	err := r.client.Del(r.ctx, token).Err()
	if err != nil {
		return fmt.Errorf("failed to delete token from Redis: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}
