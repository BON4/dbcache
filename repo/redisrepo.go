package repo

import (
	"context"
	"log"

	redis "github.com/go-redis/redis/v8"
)

type RedisCache struct {
	conn *redis.Client
}

func NewRedisCache(cache *redis.Client) Cache {
	return &RedisCache{conn: cache}
}

func (r *RedisCache) Get(ctx context.Context, key string) (*map[string]string, error) {
	redisResp, err := r.conn.HGetAll(ctx, key).Result()

	if err != nil {
		return nil, err
	}

	if len(redisResp) == 0 {
		return nil, redis.Nil
	}

	return &redisResp, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, val map[string]interface{}) error {
	err := r.conn.HSet(ctx, key, val).Err()
	if err != nil {
		return err
	}
	log.Println("Seted in redis", key, val)
	return nil
}
