package redisclient

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/yigithankarabulut/distributed-mail-queue-service/config"
	"time"
)

var Rdb *redis.Client

func GetRedisClient() *redis.Client {
	return Rdb
}

func New(config config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: "",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	Rdb = client
	log.Infof("connected to redis: %s:%s", config.Host, config.Port)
	return client, nil
}
