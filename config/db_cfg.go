package config

type mysqlConfig struct {
	Url          string `yaml:"url"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxLifetime  int64  `yaml:"maxLifetime"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxIdleTime  int64  `yaml:"maxIdleTime"`
	PrepareStmt  bool   `yaml:"prepareStmt"`
}

type dbConfig struct {
	Enable bool        `yaml:"enable"`
	Type   string      `yaml:"type"`
	Mysql  mysqlConfig `yaml:"mysql"`
}
