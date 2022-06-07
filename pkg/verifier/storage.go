package verifier

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// TokenStorage token 存储器
type TokenStorage interface {
	Get(key string) (string, error) // return error when key does not exist
	Set(key string, value interface{}, expiration time.Duration) error
	SetNX(key string, value interface{}, expiration time.Duration) bool
	Del(key string) error
	DelByKeyPrefix(keyPrefix string) error
	Exists(key string) bool
}

type RedisStorage struct {
	Cli redis.Cmdable
}

// a non-nil, empty Context
var backgroundContext = context.Background()

func (r *RedisStorage) Get(key string) (string, error) {
	return r.Cli.Get(backgroundContext, key).Result()
}

func (r *RedisStorage) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Cli.Set(backgroundContext, key, value, expiration).Err()
}

func (r *RedisStorage) SetNX(key string, value interface{}, expiration time.Duration) bool {
	return r.Cli.SetNX(backgroundContext, key, value, expiration).Val()
}

func (r *RedisStorage) Del(key string) error {
	return r.Cli.Del(backgroundContext, key).Err()
}

func (r *RedisStorage) DelByKeyPrefix(keyPrefix string) error {
	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = r.Cli.Scan(backgroundContext, cursor, keyPrefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err = r.Cli.Del(backgroundContext, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *RedisStorage) Exists(key string) bool {
	return r.Cli.Exists(backgroundContext, key).Val() == 1
}
