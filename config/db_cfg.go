package config

type mysqlConfig struct {
	Url          string `yaml:"url"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxLifetime  int64  `yaml:"maxLifetime"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxIdleTime  int64  `yaml:"maxIdleTime"`
	PrepareStmt  bool   `yaml:"prepareStmt"`
	// GORM perform single create, update, delete operations in transactions by default to ensure database data integrity
	// You can disable it by setting `SkipDefaultTransaction` to true
	SkipDefaultTransaction bool `yaml:"skipDefaultTransaction"`
	// FullSaveAssociations full save associations
	FullSaveAssociations bool
	// DryRun generate sql without execute
	DryRun bool
	// DisableAutomaticPing
	DisableAutomaticPing bool
	// DisableForeignKeyConstraintWhenMigrating
	DisableForeignKeyConstraintWhenMigrating bool
	// IgnoreRelationshipsWhenMigrating
	IgnoreRelationshipsWhenMigrating bool
	// DisableNestedTransaction disable nested transaction
	DisableNestedTransaction bool
	// AllowGlobalUpdate allow global update
	AllowGlobalUpdate bool
	// QueryFields executes the SQL query with all fields of the table
	QueryFields bool
	// CreateBatchSize default create batch size
	CreateBatchSize int
	// TranslateError enabling error translation
	TranslateError bool
}

type dbConfig struct {
	Enable bool        `yaml:"enable"`
	Type   string      `yaml:"type"`
	Mysql  mysqlConfig `yaml:"mysql"`
}
