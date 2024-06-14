package config

import (
	"github.com/caixr9527/go-cloud/component"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"math"
	"reflect"
	"sync"
)

type Configuration config

var once sync.Once

const (
	PROD = "prod"
	DEV  = "dev"
	TEST = "test"
)

type config struct {
	Server    serverConfig   `yaml:"server" mapstructure:"server"`
	Logger    logConfig      `yaml:"logger" mapstructure:"logger"`
	Cloud     cloudConfig    `yaml:"cloud" mapstructure:"cloud"`
	BasicAuth basicAuth      `yaml:"basicAuth" mapstructure:"basicAuth"`
	Jwt       jwt            `yaml:"jwt" mapstructure:"jwt"`
	Template  template       `yaml:"template" mapstructure:"template"`
	Db        dbConfig       `yaml:"db" mapstructure:"db"`
	Redis     redisConfig    `yaml:"redis" mapstructure:"redis"`
	Discover  discoverConfig `yaml:"discoverConfig" mapstructure:"discoverConfig"`
}

type conf struct {
}

func (c *conf) Create(s *component.Singleton) {
	once.Do(func() {
		var cfg Configuration
		viper.SetConfigFile("conf/application.yaml")
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Println("reload config")
			err := viper.Unmarshal(&cfg)
			if err != nil {
				log.Println(err)
			}
		})
		err := viper.ReadInConfig()
		if err != nil {
			log.Println(err)
		}
		err = viper.Unmarshal(&cfg)
		if err != nil {
			log.Println(err)
		}
		s.Register(reflect.TypeOf(cfg).String(), cfg)
	})
}

func init() {
	component.RegisterComponent(&conf{})
}

func (c *conf) Order() int {
	return math.MinInt
}
