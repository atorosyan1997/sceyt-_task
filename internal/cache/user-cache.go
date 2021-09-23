package cache

import "sceyt_task/internal/data"

type UserCache interface {
	Set(key string, value *data.User) error
	Get(key string) (*data.User, error)
	Del(key string) error
}
