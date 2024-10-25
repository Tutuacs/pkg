package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Addr string
	conn *redis.Client
}

var ctx = context.Background()

func NewRedisClient(addr string) (client *RedisClient, err error) {

	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err = conn.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	client = &RedisClient{
		Addr: addr,
		conn: conn,
	}

	return client, nil
}

func (r *RedisClient) Set(key string, value interface{}, exp time.Duration, timeout time.Duration) error {
	if timeout == 0 {
		return r.conn.Set(ctx, key, value, exp).Err()
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return r.conn.Set(ctx, key, value, exp).Err()
}

func (r *RedisClient) Get(key string, timeout time.Duration) (string, error) {
	if timeout == 0 {
		return r.conn.Get(ctx, key).Result()
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return r.conn.Get(ctx, key).Result()
}
