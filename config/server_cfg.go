package config

type https struct {
	Enable   bool   `yaml:"enable"`
	CertPath string `yaml:"cert"`
	KeyPath  string `yaml:"key"`
}

type serverConfig struct {
	Port  uint32 `yaml:"port"`
	Https https  `yaml:"https"`
}
