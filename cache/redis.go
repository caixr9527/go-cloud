package cache

import (
	"context"
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/log"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Redis struct {
	client  *redis.Client
	context context.Context
}

var RedisClient *Redis
var once sync.Once

func Init() {
	if !config.Cfg.Redis.Enable {
		return
	}
	once.Do(func() {
		log.Log.Info("init mysql conn")
		client := redis.NewClient(&redis.Options{
			Network:               config.Cfg.Redis.Network,
			Addr:                  config.Cfg.Redis.Addr,
			ClientName:            config.Cfg.Redis.ClientName,
			Protocol:              config.Cfg.Redis.Protocol,
			Username:              config.Cfg.Redis.Username,
			Password:              config.Cfg.Redis.Password,
			DB:                    config.Cfg.Redis.DB,
			MaxRetries:            config.Cfg.Redis.MaxRetries,
			MinRetryBackoff:       time.Duration(config.Cfg.Redis.MinRetryBackoff),
			MaxRetryBackoff:       time.Duration(config.Cfg.Redis.MaxRetryBackoff),
			DialTimeout:           time.Duration(config.Cfg.Redis.DialTimeout),
			ReadTimeout:           time.Duration(config.Cfg.Redis.ReadTimeout),
			WriteTimeout:          time.Duration(config.Cfg.Redis.WriteTimeout),
			ContextTimeoutEnabled: config.Cfg.Redis.ContextTimeoutEnabled,
			PoolFIFO:              config.Cfg.Redis.PoolFIFO,
			PoolSize:              config.Cfg.Redis.PoolSize,
			PoolTimeout:           time.Duration(config.Cfg.Redis.PoolTimeout),
			MinIdleConns:          config.Cfg.Redis.MinIdleConns,
			MaxIdleConns:          config.Cfg.Redis.MaxIdleConns,
			MaxActiveConns:        config.Cfg.Redis.MaxActiveConns,
			ConnMaxIdleTime:       time.Duration(config.Cfg.Redis.ConnMaxIdleTime),
			ConnMaxLifetime:       time.Duration(config.Cfg.Redis.ConnMaxLifetime),
			DisableIndentity:      config.Cfg.Redis.DisableIndentity,
			IdentitySuffix:        config.Cfg.Redis.IdentitySuffix,
		})
		RedisClient = &Redis{
			client:  client,
			context: context.Background(),
		}
		log.Log.Info("init mysql conn success")
	})
}

func (r *Redis) Set(key string, value any) {
	err := r.client.Set(r.context, key, value, 0).Err()
	r.printLog(err)
}

func (r *Redis) SetEx(key string, value any, expire time.Duration) {
	err := r.client.SetEx(r.context, key, value, expire).Err()
	r.printLog(err)
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(r.context, key).Result()
}

func (r *Redis) HSet(key string, value ...any) {
	err := r.client.HSet(r.context, key, value...).Err()
	r.printLog(err)
}

func (r *Redis) HMSet(key string, value ...any) {
	err := r.client.HMSet(r.context, key, value...).Err()
	r.printLog(err)
}

func (r *Redis) HGet(key string, field string) (string, error) {
	return r.client.HGet(r.context, key, field).Result()
}

func (r *Redis) HMGet(key string, field ...string) ([]any, error) {
	return r.client.HMGet(r.context, key, field...).Result()
}

func (r *Redis) LPush(key string, value ...any) {
	err := r.client.LPush(r.context, key, value...).Err()
	r.printLog(err)
}

func (r *Redis) RPush(key string, value ...any) {
	err := r.client.RPush(r.context, key, value...).Err()
	r.printLog(err)
}

func (r *Redis) LPop(key string) (string, error) {
	return r.client.LPop(r.context, key).Result()
}

func (r *Redis) RPop(key string) (string, error) {
	return r.client.RPop(r.context, key).Result()
}

func (r *Redis) Incr(key string, value ...int64) (int64, error) {
	if len(value) > 0 {
		return r.client.IncrBy(r.context, key, value[0]).Result()
	} else {
		return r.client.Incr(r.context, key).Result()
	}
}

func (r *Redis) Decr(key string, value ...int64) (int64, error) {
	if len(value) > 0 {
		return r.client.DecrBy(r.context, key, value[0]).Result()
	} else {
		return r.client.Decr(r.context, key).Result()
	}
}

func (r *Redis) Del(key string) {
	err := r.client.Del(r.context, key).Err()
	r.printLog(err)
}

func (r *Redis) Expire(key string, expiration time.Duration) {
	err := r.client.Expire(r.context, key, expiration).Err()
	r.printLog(err)
}

func (r *Redis) Exists(keys ...string) (int64, error) {
	return r.client.Exists(r.context, keys...).Result()
}

func (r *Redis) printLog(err error) {
	if err != nil {
		log.Log.Error(err.Error())
	}
}
