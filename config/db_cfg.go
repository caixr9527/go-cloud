package config

type mysqlConfig struct {
	Url          string `yaml:"url" mapstructure:"url"`
	MaxOpenConns int    `yaml:"maxOpenConns" mapstructure:"maxOpenConns"`
	// Millisecond
	MaxLifetime  int64 `yaml:"maxLifetime" mapstructure:"maxLifetime"`
	MaxIdleConns int   `yaml:"maxIdleConns" mapstructure:"maxIdleConns"`
	// Millisecond
	MaxIdleTime int64 `yaml:"maxIdleTime" mapstructure:"maxIdleTime"`
	PrepareStmt bool  `yaml:"prepareStmt" mapstructure:"prepareStmt"`
	// GORM perform single create, update, delete operations in transactions by default to ensure database data integrity
	// You can disable it by setting `SkipDefaultTransaction` to true
	SkipDefaultTransaction bool `yaml:"skipDefaultTransaction" mapstructure:"skipDefaultTransaction"`
	// FullSaveAssociations full save associations
	FullSaveAssociations bool `yaml:"fullSaveAssociations" mapstructure:"fullSaveAssociations"`
	// DryRun generate sql without execute
	DryRun bool `yaml:"dryRun" mapstructure:"dryRun"`
	// DisableAutomaticPing
	DisableAutomaticPing bool `yaml:"disableAutomaticPing" mapstructure:"disableAutomaticPing"`
	// DisableForeignKeyConstraintWhenMigrating
	DisableForeignKeyConstraintWhenMigrating bool `yaml:"disableForeignKeyConstraintWhenMigrating" mapstructure:"disableForeignKeyConstraintWhenMigrating"`
	// IgnoreRelationshipsWhenMigrating
	IgnoreRelationshipsWhenMigrating bool `yaml:"ignoreRelationshipsWhenMigrating" mapstructure:"ignoreRelationshipsWhenMigrating"`
	// DisableNestedTransaction disable nested transaction
	DisableNestedTransaction bool `yaml:"disableNestedTransaction" mapstructure:"disableNestedTransaction"`
	// AllowGlobalUpdate allow global update
	AllowGlobalUpdate bool `yaml:"allowGlobalUpdate" mapstructure:"allowGlobalUpdate"`
	// QueryFields executes the SQL query with all fields of the table
	QueryFields bool `yaml:"queryFields" mapstructure:"queryFields"`
	// CreateBatchSize default create batch size
	CreateBatchSize int `yaml:"createBatchSize" mapstructure:"createBatchSize"`
	// TranslateError enabling error translation
	TranslateError bool `yaml:"translateError" mapstructure:"translateError"`
}

type dbConfig struct {
	Enable bool        `yaml:"enable" mapstructure:"enable"`
	Type   string      `yaml:"type" mapstructure:"type"`
	Mysql  mysqlConfig `yaml:"mysql" mapstructure:"mysql"`
}
