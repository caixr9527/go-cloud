package config

import (
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"math"
	"reflect"
	"strings"
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

//type Configuration struct {
//}

func (c *Configuration) Create(s *component.Singleton) {
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

func (c *Configuration) Refresh(s *component.Singleton) {
	configClient := factory.Get(&config_client.ConfigClient{})
	configuration := factory.Get(Configuration{})
	dataIds := strings.Split(configuration.Discover.Config.DataIds, ",")
	group := configuration.Discover.Config.Group
	var contents strings.Builder
	for index := range dataIds {
		dataId := dataIds[index]
		content, err := configClient.GetConfig(vo.ConfigParam{
			DataId: dataId,
			Group:  group,
		})
		if err != nil {
			log.Println(err)
		} else {
			contents.WriteString(content)
		}
	}

	err := yaml.Unmarshal([]byte(contents.String()), &configuration)
	if err != nil {
		log.Println(err)
	}
	s.Register(reflect.TypeOf(configuration).String(), configuration)

}

func init() {
	component.RegisterComponent(&Configuration{})
}

func (c *Configuration) Order() int {
	return math.MinInt
}
