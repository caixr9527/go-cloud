package cache

import (
	"context"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/log"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	Client  *redis.Client
	Context context.Context
}

func (r *Redis) Set(key string, value any) {
	err := r.Client.Set(r.Context, key, value, 0).Err()
	r.printLog(err)
}

func (r *Redis) SetEx(key string, value any, expire time.Duration) {
	err := r.Client.SetEx(r.Context, key, value, expire).Err()
	r.printLog(err)
}

func (r *Redis) Get(key string) (string, error) {
	return r.Client.Get(r.Context, key).Result()
}

func (r *Redis) HSet(key string, value ...any) {
	err := r.Client.HSet(r.Context, key, value...).Err()
	r.printLog(err)
}

func (r *Redis) HMSet(key string, value ...any) {
	err := r.Client.HMSet(r.Context, key, value...).Err()
	r.printLog(err)
}

func (r *Redis) HGet(key string, field string) (string, error) {
	return r.Client.HGet(r.Context, key, field).Result()
}

func (r *Redis) HMGet(key string, field ...string) ([]any, error) {
	return r.Client.HMGet(r.Context, key, field...).Result()
}

func (r *Redis) LPush(key string, value ...any) {
	err := r.Client.LPush(r.Context, key, value...).Err()
	r.printLog(err)
}

func (r *Redis) RPush(key string, value ...any) {
	err := r.Client.RPush(r.Context, key, value...).Err()
	r.printLog(err)
}

func (r *Redis) LPop(key string) (string, error) {
	return r.Client.LPop(r.Context, key).Result()
}

func (r *Redis) RPop(key string) (string, error) {
	return r.Client.RPop(r.Context, key).Result()
}

func (r *Redis) Incr(key string, value ...int64) (int64, error) {
	if len(value) > 0 {
		return r.Client.IncrBy(r.Context, key, value[0]).Result()
	} else {
		return r.Client.Incr(r.Context, key).Result()
	}
}

func (r *Redis) Decr(key string, value ...int64) (int64, error) {
	if len(value) > 0 {
		return r.Client.DecrBy(r.Context, key, value[0]).Result()
	} else {
		return r.Client.Decr(r.Context, key).Result()
	}
}

func (r *Redis) Del(key string) {
	err := r.Client.Del(r.Context, key).Err()
	r.printLog(err)
}

func (r *Redis) Expire(key string, expiration time.Duration) {
	err := r.Client.Expire(r.Context, key, expiration).Err()
	r.printLog(err)
}

func (r *Redis) Exists(keys ...string) (int64, error) {
	return r.Client.Exists(r.Context, keys...).Result()
}

func (r *Redis) printLog(err error) {
	if err != nil {
		logger := factory.Get(&log.Log{})
		logger.Error(err.Error())
	}
}
