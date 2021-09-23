package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"sceyt_task/internal/data"
	"sceyt_task/pkg/logging"
	"time"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
	logger  logging.Logger
}

func NewRedisCache(host string, db int, exp time.Duration) UserCache {
	return &redisCache{host: host, db: db, expires: exp}
}

func (r *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     r.host,
		Password: "",
		DB:       r.db,
	})
}

func (r *redisCache) Set(key string, value *data.User) error {
	client := r.getClient()
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	client.Set(key, json, r.expires*time.Second)
	_, err = client.Get(key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (r *redisCache) Get(key string) (*data.User, error) {
	client := r.getClient()
	res, err := client.Get(key).Result()
	if err != nil {
		return nil, nil
	}
	user := &data.User{}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *redisCache) Del(key string) error {
	client := r.getClient()
	_, err := client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}
