package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
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
	Server    serverConfig `yaml:"server"`
	Logger    logConfig    `yaml:"logger"`
	Cloud     cloudConfig  `yaml:"cloud"`
	BasicAuth basicAuth    `yaml:"basicAuth"`
	Jwt       jwt          `yaml:"jwt"`
	Template  template     `yaml:"template"`
	Db        dbConfig     `yaml:"db"`
}

func Init() {
	once.Do(func() {
		data := loadYaml()
		loadConfig(data)
	})
}

func loadConfig(data []byte) {
	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		log.Fatal(err)
	}
}

func loadYaml() []byte {
	data, err := os.ReadFile("conf/application.yaml")
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
		return nil
	}
	return data
}
