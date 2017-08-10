package redis

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisClient struct {
	Addr     string
	Password string
	DB       int
	Client   *redis.Client
}

func (rc *RedisClient) Connect() {
	if rc.Client != nil {
		pong, err := rc.Client.Ping().Result()
		if err == nil && pong == "PONG" {
			return
		}
	}
	rc.Client = redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
	})
}

func (rc *RedisClient) GetKeys(pattern string) ([]string, error) {
	return rc.Client.Keys(pattern).Result()
}

func (rc *RedisClient) SetString(key, val string) error {
	err := rc.Client.Set(key, val, 0).Err()
	return err
}

func (rc *RedisClient) GetString(key string) (string, error) {
	return rc.Client.Get(key).Result()
}

func (rc *RedisClient) SetHash(key string, fields map[string]interface{}) error {
	return rc.Client.HMSet(key, fields).Err()
}

func (rc *RedisClient) GetHash(key string) (map[string]string, error) {
	return rc.Client.HGetAll(key).Result()
}

func (rc *RedisClient) SetList(key, val string) error {
	return rc.Client.LPush(key, val).Err()
}

func (rc *RedisClient) GetList(key string, start, stop int64) ([]string, error) {
	return rc.Client.LRange(key, start, stop).Result()
}

func (rc *RedisClient) SetExpires(key string, expiration time.Duration) error {
	return rc.Client.Expire(key, expiration).Err()
}

func (rc *RedisClient) ExpiresAt(key string, t time.Time) error {
	return rc.Client.ExpireAt(key, t).Err()
}
