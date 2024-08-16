package cache

import (
	"context"
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/factory"
	"github.com/caixr9527/go-cloud/log"
	"github.com/redis/go-redis/v9"
	"math"
	"sync"
	"time"
)

var once sync.Once

type Cache struct {
	Redis *Redis
}

func init() {
	factory.RegisterComponent(&Cache{})
}

func (c *Cache) Order() int {
	return math.MinInt + 3
}

func (c *Cache) Destroy() {
	cache := factory.Get(c)
	cache.Redis.Client.Close()
	factory.Del(c)
	factory.Get(&log.Log{}).Info("redis destroy success")
}
func (c *Cache) Refresh() {

}

func (c *Cache) Name() string {
	return utils.ObjName(c)
}

func (c *Cache) Create() {
	configuration := factory.Get(&config.Configuration{})
	if !configuration.Redis.Enable {
		return
	}
	logger := factory.Get(&log.Log{})
	once.Do(func() {
		logger.Info("create redis conn")
		c.new(configuration)
		logger.Info("create redis conn success")
	})
}

func (c *Cache) new(configuration *config.Configuration) {
	client := redis.NewClient(&redis.Options{
		Network:               configuration.Redis.Network,
		Addr:                  configuration.Redis.Addr,
		ClientName:            configuration.Redis.ClientName,
		Protocol:              configuration.Redis.Protocol,
		Username:              configuration.Redis.Username,
		Password:              configuration.Redis.Password,
		DB:                    configuration.Redis.DB,
		MaxRetries:            configuration.Redis.MaxRetries,
		MinRetryBackoff:       time.Duration(configuration.Redis.MinRetryBackoff),
		MaxRetryBackoff:       time.Duration(configuration.Redis.MaxRetryBackoff),
		DialTimeout:           time.Duration(configuration.Redis.DialTimeout),
		ReadTimeout:           time.Duration(configuration.Redis.ReadTimeout),
		WriteTimeout:          time.Duration(configuration.Redis.WriteTimeout),
		ContextTimeoutEnabled: configuration.Redis.ContextTimeoutEnabled,
		PoolFIFO:              configuration.Redis.PoolFIFO,
		PoolSize:              configuration.Redis.PoolSize,
		PoolTimeout:           time.Duration(configuration.Redis.PoolTimeout),
		MinIdleConns:          configuration.Redis.MinIdleConns,
		MaxIdleConns:          configuration.Redis.MaxIdleConns,
		MaxActiveConns:        configuration.Redis.MaxActiveConns,
		ConnMaxIdleTime:       time.Duration(configuration.Redis.ConnMaxIdleTime),
		ConnMaxLifetime:       time.Duration(configuration.Redis.ConnMaxLifetime),
		DisableIndentity:      configuration.Redis.DisableIndentity,
		IdentitySuffix:        configuration.Redis.IdentitySuffix,
	})
	redisClient := &Redis{
		Client:  client,
		Context: context.Background(),
	}
	c.Redis = redisClient
	factory.Create(c)
}
