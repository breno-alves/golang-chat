package cache

import (
	"chatroom/config"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

func NewCache(config *config.Config) *redis.Client {
	defer func() {
		slog.Info("successfully connected to redis")
	}()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Cache.Host,
		Password: config.Cache.Password,
		DB:       config.Cache.Db,
	})
	return rdb
}
