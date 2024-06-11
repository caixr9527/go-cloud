package config

type https struct {
	Enable   bool   `yaml:"enable" mapstructure:"enable"`
	CertPath string `yaml:"cert" mapstructure:"cert"`
	KeyPath  string `yaml:"key" mapstructure:"key"`
}

type serverConfig struct {
	Port  uint32 `yaml:"port" mapstructure:"port"`
	Https https  `yaml:"https" mapstructure:"https"`
}
