package GorraAPI

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func InitRedis(host string, port int, db int, password string) *redis.Client {
	cfg := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
		DB:   db,
	}

	if password != "" {
		cfg.Password = password
	}

	rdb := redis.NewClient(cfg)

	if rdb.Ping(context.Background()).Err() != nil {
		zap.S().Panicw("redis init failed", "err", "redis init failed")
	}

	return rdb
}
