package config

import (
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"math"
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

func (c *Configuration) Name() string {
	return utils.ObjName(c)
}
func (c *Configuration) Create() {
	once.Do(func() {
		//var cfg Configuration
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
	})
}

func (c *Configuration) Destroy() {
	factory.Del(c)
	factory.Get(&zap.Logger{}).Info("Configuration destroy success")
}

func (c *Configuration) Refresh() {
	// todo
	//configClient := factory.GetField("discover.Discover", &config_client.ConfigClient{})
	//configuration := factory.Get(&Configuration{})
	//dataIds := strings.Split(configuration.Discover.Config.DataIds, ",")
	//group := configuration.Discover.Config.Group
	//var contents strings.Builder
	//for index := range dataIds {
	//	dataId := dataIds[index]
	//	content, err := configClient.GetConfig(vo.ConfigParam{
	//		DataId: dataId,
	//		Group:  group,
	//	})
	//	if err != nil {
	//		log.Println(err)
	//	} else {
	//		contents.WriteString(content)
	//	}
	//}
	//
	//err := yaml.Unmarshal([]byte(contents.String()), &configuration)
	//if err != nil {
	//	log.Println(err)
	//}
	//factory.Create(configuration)

}

func init() {
	component.RegisterComponent(&Configuration{})
}

func (c *Configuration) Order() int {
	return math.MinInt
}
