package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	conf "github.com/birukbelay/gocmn/src/config"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/provider/db"
)

type RedisService struct {
	RedisClient *redis.Client
}

// Exists implements providers.KeyValServ.
func NewRedis(config *conf.KeyValConfig) (*RedisService, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.KVHost, config.KVPort),
		Password: config.KVPassword, // no password set if empty
		DB:       config.KVDbName,   // use default DB
		Username: config.KVUsername,
	})
	if rdb != nil {
		pong, rErr := rdb.Ping(context.Background()).Result()
		if rErr != nil {
			logger.LogError("redis connection error", rErr.Error())
			return nil, rErr
		}
		logger.LogInfo("Redis success", pong)
		return &RedisService{RedisClient: rdb}, nil
	}
	return nil, fmt.Errorf("redis client is nil")
}

// Exists implements providers.KeyValServ.
func NewKeyValServ(config *conf.KeyValConfig) (db.KeyValServ, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.KVHost, config.KVPort),
		Password: config.KVPassword, // no password set if empty
		DB:       config.KVDbName,   // use default DB
		Username: config.KVUsername,
	})
	return &RedisService{RedisClient: rdb}, nil
}

func (r *RedisService) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if exists > 0 {
		return true, nil
	}
	return false, nil
}

// Get implements providers.KeyValServ.
func (r *RedisService) Get(ctx context.Context, key string) (value any, exists bool, er error) {
	value, err := r.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return value, false, nil
	}
	if err != nil {
		return value, false, err
	}
	return value, true, nil
}

// Set implements providers.KeyValServ.
func (r *RedisService) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := r.RedisClient.Set(ctx, key, val, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value: %v", err)
	}
	return nil
}
