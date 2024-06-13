package cache

import (
	"context"
	redis2 "github.com/caixr9527/go-cloud/cache/redis"
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/log"
	"github.com/redis/go-redis/v9"
	"math"
	"sync"
	"time"
)

var once sync.Once

type cache struct {
}

func init() {
	component.RegisterComponent(&cache{})
}

func (c *cache) Order() int {
	return math.MinInt + 3
}

func (c *cache) Create(s *component.Singleton) {
	if !factory.GetConf().Redis.Enable {
		return
	}
	once.Do(func() {
		log.Log.Info("init redis conn")
		client := redis.NewClient(&redis.Options{
			Network:               factory.GetConf().Redis.Network,
			Addr:                  factory.GetConf().Redis.Addr,
			ClientName:            factory.GetConf().Redis.ClientName,
			Protocol:              factory.GetConf().Redis.Protocol,
			Username:              factory.GetConf().Redis.Username,
			Password:              factory.GetConf().Redis.Password,
			DB:                    factory.GetConf().Redis.DB,
			MaxRetries:            factory.GetConf().Redis.MaxRetries,
			MinRetryBackoff:       time.Duration(factory.GetConf().Redis.MinRetryBackoff),
			MaxRetryBackoff:       time.Duration(factory.GetConf().Redis.MaxRetryBackoff),
			DialTimeout:           time.Duration(factory.GetConf().Redis.DialTimeout),
			ReadTimeout:           time.Duration(factory.GetConf().Redis.ReadTimeout),
			WriteTimeout:          time.Duration(factory.GetConf().Redis.WriteTimeout),
			ContextTimeoutEnabled: factory.GetConf().Redis.ContextTimeoutEnabled,
			PoolFIFO:              factory.GetConf().Redis.PoolFIFO,
			PoolSize:              factory.GetConf().Redis.PoolSize,
			PoolTimeout:           time.Duration(factory.GetConf().Redis.PoolTimeout),
			MinIdleConns:          factory.GetConf().Redis.MinIdleConns,
			MaxIdleConns:          factory.GetConf().Redis.MaxIdleConns,
			MaxActiveConns:        factory.GetConf().Redis.MaxActiveConns,
			ConnMaxIdleTime:       time.Duration(factory.GetConf().Redis.ConnMaxIdleTime),
			ConnMaxLifetime:       time.Duration(factory.GetConf().Redis.ConnMaxLifetime),
			DisableIndentity:      factory.GetConf().Redis.DisableIndentity,
			IdentitySuffix:        factory.GetConf().Redis.IdentitySuffix,
		})
		redisClient := &redis2.Redis{
			Client:  client,
			Context: context.Background(),
		}
		s.Register("redis", redisClient)
		log.Log.Info("init redis conn success")
	})
}
