package cache

import (
	"context"
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math"
	"reflect"
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

func (c *cache) Refresh(s *component.Singleton) {
	configuration := factory.Get(config.Configuration{})
	if !configuration.Redis.Enable {
		return
	}
	logger := factory.Get(&zap.Logger{})
	logger.Info("refresh redis conn")
	c.new(configuration, s)
	logger.Info("refresh redis conn success")
}

func (c *cache) Create(s *component.Singleton) {
	configuration := factory.Get(config.Configuration{})
	if !configuration.Redis.Enable {
		return
	}
	logger := factory.Get(&zap.Logger{})
	once.Do(func() {
		logger.Info("create redis conn")
		c.new(configuration, s)
		logger.Info("create redis conn success")
	})
}

func (c *cache) new(configuration config.Configuration, s *component.Singleton) {
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
	s.Register(reflect.TypeOf(redisClient).String(), redisClient)
}
