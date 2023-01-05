package entity

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisEntity interface {
	Set(key string, value interface{}, expire time.Duration) (string, error)
	Get(key string) (interface{}, error)
	Del(key string) (interface{}, error)
}

type redisConnection struct {
	connection *redis.Client
}

func NewRedisEntity(rdb *redis.Client) RedisEntity {
	return &redisConnection{
		connection: rdb,
	}
}

func (rdb *redisConnection) Set(key string, value interface{}, expire time.Duration) (string, error) {
	val, err := rdb.connection.Set(ctx, key, value, expire).Result()
	return val, err
}

func (rdb *redisConnection) Get(key string) (interface{}, error) {
	val, err := rdb.connection.Get(ctx, key).Result()
	return val, err
}

func (rdb *redisConnection) Del(key string) (interface{}, error) {
	val, err := rdb.connection.Del(ctx, key).Result()
	return val, err
}
