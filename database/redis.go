package database

import (
	"context"
	"fmt"
	"log"
	"marketplace-system/config"

	"github.com/go-redis/redis/v8"
)

func GetConnectionRedis(redisCred config.Redis) (rdb *redis.Client, err error) {
	ctx := context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%s`, redisCred.Host, redisCred.Port),
		Password: redisCred.Password,
		DB:       redisCred.Db,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	return
}
