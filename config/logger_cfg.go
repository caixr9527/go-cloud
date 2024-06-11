package config

type logConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`
	FileName   string `yaml:"fileName" mapstructure:"fileName"`
	MaxSize    uint32 `yaml:"maxSize" mapstructure:"maxSize"`
	MaxAge     uint32 `yaml:"maxAge" mapstructure:"maxAge"`
	MaxBackups uint32 `yaml:"maxBackups" mapstructure:"maxBackups"`
	Compress   bool   `yaml:"compress" mapstructure:"compress"`
}
