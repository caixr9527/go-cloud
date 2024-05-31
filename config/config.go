package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var Cfg config
var once sync.Once

type config struct {
	Server serverConfig `yaml:"server"`
	Logger logConfig    `yaml:"logger"`
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
