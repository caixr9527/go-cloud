package config

import (
	"errors"
	"fmt"
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/factory"
	"github.com/caixr9527/go-cloud/internal/component"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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
	Discover  discoverConfig `yaml:"discover" mapstructure:"discover"`
}

var configs = make([]any, 0)

func Add(conf ...any) error {
	if len(conf) == 0 {
		return errors.New("invalid parameters")
	}
	for index := range conf {
		c := conf[index]
		typeOf := reflect.TypeOf(c)
		if typeOf.Kind() != reflect.Pointer {
			return errors.New(typeOf.String() + " must be pointer")
		}
	}
	configs = append(configs, conf...)
	return nil
}

func (c *Configuration) Name() string {
	return utils.ObjName(c)
}
func (c *Configuration) Create() {
	once.Do(func() {
		viper.SetConfigFile("conf/application.yaml")
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Println("reload config")
			err := viper.Unmarshal(&c)
			if err != nil {
				log.Println(err)
			}
		})
		err := viper.ReadInConfig()
		if err != nil {
			log.Println(err)
		}
		err = viper.Unmarshal(&c)
		if err != nil {
			log.Println(err)
		}
		factory.Create(c)
		c.LoadLocalCustomConfig()
	})
}

func (c *Configuration) LoadLocalCustomConfig() {
	for idx := range configs {
		conf := configs[idx]
		err := viper.Unmarshal(conf)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Configuration) LoadRemoteCustomConfig(data string) {
	for idx := range configs {
		conf := configs[idx]
		err := yaml.Unmarshal([]byte(data), conf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(fmt.Sprintf("refresh config: %v success", conf))
	}
}

func (c *Configuration) Destroy() {
	factory.Del(c)
	log.Println("destroy configuration success")
}

func (c *Configuration) Refresh() {

}

func init() {
	component.RegisterComponent(&Configuration{})
}

func (c *Configuration) Order() int {
	return math.MinInt
}
