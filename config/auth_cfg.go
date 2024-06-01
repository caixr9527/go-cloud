package config

import "time"

type basicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type jwt struct {
	Header       string        `yaml:"header"`
	SecretKey    string        `yaml:"secretKey"`
	Alg          string        `yaml:"alg"`
	TokenTimeout time.Duration `yaml:"timeout"`
	RefreshKey   string        `yaml:"refreshKey"`
	Whitelist    []string      `yaml:"ignore"`
}
