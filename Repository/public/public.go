package public

import (
	"encoding/json"
	"fmt"
	"kaktus-task/Connections"
	"time"
)

func (p public) RedisSetCache(keys string, data interface{}, expired time.Duration) (err error) {
	conn := Connections.RedisConn
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = conn.Set(keys, string(dataBytes), expired).Err()
	return
}

func (p public) RedisGetCache(keys string, data interface{}) (err error) {
	conn := Connections.RedisConn
	dataCache, err := conn.Get(keys).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(dataCache), &data)
	return
}

func (p public) RedisSetHMCache(keys string, data map[string]interface{}, expired time.Duration) (err error) {
	conn := Connections.RedisConn
	err = conn.HMSet(keys, data).Err()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = conn.Expire(keys, expired).Err()
	return
}

func (p public) RedisGetHMCache(keys, field string, data interface{}) (err error) {
	conn := Connections.RedisConn
	dataRedis, err := conn.HMGet(keys, field).Result()
	if err != nil {
		return
	}
	if len(dataRedis) == 0 || dataRedis[0] == nil {
		return
	}
	err = json.Unmarshal([]byte(dataRedis[0].(string)), &data)
	return
}

func (p public) RedisHGetAllCache(keys string) (res map[string]string, err error) {
	conn := Connections.RedisConn
	res, err = conn.HGetAll(keys).Result()
	return
}

func (p public) RedisHDelCache(keys string, field ...string) (err error) {
	conn := Connections.RedisConn
	err = conn.HDel(keys, field...).Err()
	return
}

func (p public) RedisDelCache(keys ...string) (err error) {
	conn := Connections.RedisConn
	err = conn.Del(keys...).Err()
	return
}
