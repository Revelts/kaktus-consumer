package public

import (
	"time"
)

type PublicInterface interface {
	RedisSetCache(keys string, data interface{}, expired time.Duration) (err error)
	RedisGetCache(keys string, data interface{}) (err error)
	RedisSetHMCache(keys string, data map[string]interface{}, expired time.Duration) (err error)
	RedisGetHMCache(keys, field string, data interface{}) (err error)
	RedisHGetAllCache(keys string) (res map[string]string, err error)
	RedisHDelCache(keys string, field ...string) (err error)
	RedisDelCache(keys ...string) (err error)
}

type public struct{}

func InitPublic() PublicInterface {
	return &public{}
}
