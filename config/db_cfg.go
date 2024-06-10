package config

type mysqlConfig struct {
	Url          string `yaml:"url"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	// Millisecond
	MaxLifetime  int64 `yaml:"maxLifetime"`
	MaxIdleConns int   `yaml:"maxIdleConns"`
	// Millisecond
	MaxIdleTime int64 `yaml:"maxIdleTime"`
	PrepareStmt bool  `yaml:"prepareStmt"`
	// GORM perform single create, update, delete operations in transactions by default to ensure database data integrity
	// You can disable it by setting `SkipDefaultTransaction` to true
	SkipDefaultTransaction bool `yaml:"skipDefaultTransaction"`
	// FullSaveAssociations full save associations
	FullSaveAssociations bool `yaml:"fullSaveAssociations"`
	// DryRun generate sql without execute
	DryRun bool `yaml:"dryRun"`
	// DisableAutomaticPing
	DisableAutomaticPing bool `yaml:"disableAutomaticPing"`
	// DisableForeignKeyConstraintWhenMigrating
	DisableForeignKeyConstraintWhenMigrating bool `yaml:"disableForeignKeyConstraintWhenMigrating"`
	// IgnoreRelationshipsWhenMigrating
	IgnoreRelationshipsWhenMigrating bool `yaml:"ignoreRelationshipsWhenMigrating"`
	// DisableNestedTransaction disable nested transaction
	DisableNestedTransaction bool `yaml:"disableNestedTransaction"`
	// AllowGlobalUpdate allow global update
	AllowGlobalUpdate bool `yaml:"allowGlobalUpdate"`
	// QueryFields executes the SQL query with all fields of the table
	QueryFields bool `yaml:"queryFields"`
	// CreateBatchSize default create batch size
	CreateBatchSize int `yaml:"createBatchSize"`
	// TranslateError enabling error translation
	TranslateError bool `yaml:"translateError"`
}

type dbConfig struct {
	Enable bool        `yaml:"enable"`
	Type   string      `yaml:"type"`
	Mysql  mysqlConfig `yaml:"mysql"`
}
