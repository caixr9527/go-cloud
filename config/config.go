package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"sync"
)

var Cfg config
var once sync.Once

const (
	PROD = "prod"
	DEV  = "dev"
	TEST = "test"
)

type config struct {
	Server    serverConfig `yaml:"server" mapstructure:"server"`
	Logger    logConfig    `yaml:"logger" mapstructure:"logger"`
	Cloud     cloudConfig  `yaml:"cloud" mapstructure:"cloud"`
	BasicAuth basicAuth    `yaml:"basicAuth" mapstructure:"basicAuth"`
	Jwt       jwt          `yaml:"jwt" mapstructure:"jwt"`
	Template  template     `yaml:"template" mapstructure:"template"`
	Db        dbConfig     `yaml:"db" mapstructure:"db"`
	Redis     redisConfig  `yaml:"redis" mapstructure:"redis"`
}

func Init() {
	once.Do(func() {
		//data := loadYaml()
		//loadConfig(data)
		viper.SetConfigFile("conf/application.yaml")

		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Println("reload config")
			err := viper.Unmarshal(&Cfg)
			if err != nil {
				log.Println(err)
			}
		})
		viper.WatchConfig()
		err := viper.ReadInConfig()
		if err != nil {
			log.Println(err)
		}
		err = viper.Unmarshal(&Cfg)
		if err != nil {
			log.Println(err)
		}
	})
}
