package config

type logConfig struct {
	Level      string `yaml:"level"`
	FileName   string `yaml:"fileName"`
	MaxSize    uint32 `yaml:"maxSize"`
	MaxAge     uint32 `yaml:"maxAge"`
	MaxBackups uint32 `yaml:"maxBackups"`
	Compress   bool   `yaml:"compress"`
}
