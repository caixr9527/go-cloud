package config

import (
	"github.com/caixr9527/go-cloud/log"
	"gopkg.in/yaml.v3"
	"os"
)

var Cfg config

type config struct {
	Server serverConfig `yaml:"server"`
}

func init() {
	data := loadYaml()
	loadConfig(data)
}

func loadConfig(data []byte) {
	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		log.Log.Error(err.Error())
	}
}

func loadYaml() []byte {
	data, err := os.ReadFile("conf/application.yaml")
	if err != nil {
		log.Log.Error(err.Error())
		os.Exit(-1)
		return nil
	}
	return data
}
