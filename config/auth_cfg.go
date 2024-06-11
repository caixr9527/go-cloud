package config

import "time"

type basicAuth struct {
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	Realm    string `yaml:"realm" mapstructure:"realm"`
}

type jwt struct {
	Header       string        `yaml:"header" mapstructure:"header"`
	SecretKey    string        `yaml:"secretKey" mapstructure:"secretKey"`
	Alg          string        `yaml:"alg" mapstructure:"alg"`
	TokenTimeout time.Duration `yaml:"timeout" mapstructure:"timeout"`
	RefreshKey   string        `yaml:"refreshKey" mapstructure:"refreshKey"`
	Allow        []string      `yaml:"allow" mapstructure:"allow"`
}
