package gormkit

type Config struct {
	SkipDefaultTransaction                   bool `json:"skip_default_transaction"`
	DryRun                                   bool `json:"dry_run"`
	PrepareStmt                              bool `json:"prepare_stmt"`
	DisableAutomaticPing                     bool `json:"disable_automatic_ping"`
	DisableForeignKeyConstraintWhenMigrating bool `json:"disable_foreign_key_constraint_when_migrating"`
	IgnoreRelationshipsWhenMigrating         bool `json:"ignore_relationships_when_migrating"`
	DisableNestedTransaction                 bool `json:"disable_nested_transaction"`
}

type LoggerConfig struct {
	SlowThreshold             int64  `json:"slow_threshold"`
	IgnoreRecordNotFoundError bool   `json:"ignore_record_not_found_error"`
	ParameterizedQueries      bool   `json:"parameterized_queries"`
	LogLevel                  string `json:"log_level"`
}

type MysqlSpecificConfig struct {
	DefaultStringSize         uint `json:"default_string_size"`
	DisableDatetimePrecision  bool `json:"disable_datetime_precision"`
	DefaultDatetimePrecision  int  `json:"default_datetime_precision"`
	SupportRenameIndex        bool `json:"support_rename_index"`
	SupportRenameColumn       bool `json:"support_rename_column"`
	SkipInitializeWithVersion bool `json:"skip_initialize_with_version"`
	DisableWithReturning      bool `json:"disable_with_returning"`
	SupportForShareClause     bool `json:"support_for_share_clause"`
	SupportNullAsDefaultValue bool `json:"support_null_as_default_value"`
	SupportRenameColumnUnique bool `json:"support_rename_column_unique"`
}

type PostgresqlSpecificConfig struct {
	PreferSimpleProtocol bool  `json:"prefer_simple_protocol"`
	WithoutReturning     bool  `json:"without_returning"`
	MaxIdleConnCount     int64 `json:"max_idle_conn_count"`
	MaxOpenConnCount     int64 `json:"max_open_conn_count"`
	ConnMaxLifetime      int64 `json:"conn_max_lifetime"`
}

type Sqlite struct {
	FileName     string            `json:"db"`
	Options      map[string]string `json:"options"`
	Config       *Config           `json:"config"`
	LoggerConfig *LoggerConfig     `json:"logger"`
}

type Mysql struct {
	DatabaseName   string               `json:"db"`
	Username       string               `json:"username"`
	Password       string               `json:"password"`
	Host           string               `json:"host"`
	Port           string               `json:"port"`
	Protocol       string               `json:"protocol"`
	Options        map[string]string    `json:"options"`
	Config         *Config              `json:"config"`
	LoggerConfig   *LoggerConfig        `json:"logger"`
	SpecificConfig *MysqlSpecificConfig `json:"specific_config"`
}

type Postgresql struct {
	DatabaseName   string                    `json:"db"`
	Username       string                    `json:"username"`
	Password       string                    `json:"password"`
	Host           string                    `json:"host"`
	Port           string                    `json:"port"`
	Options        map[string]string         `json:"options"`
	Config         *Config                   `json:"config"`
	LoggerConfig   *LoggerConfig             `json:"logger"`
	SpecificConfig *PostgresqlSpecificConfig `json:"specific_config"`
}

type SqlConfigurable interface {
	Sqlite | Mysql | Postgresql
}

type MongoLoggerConfig struct {
	MaxDocumentLength   int    `json:"max_document_length"`
	ComponentCommand    string `json:"component_command"`
	ComponentConnection string `json:"component_connection"`
}

type Mongo struct {
	DatabaseName string             `json:"db"`
	Username     string             `json:"username"`
	Password     string             `json:"password"`
	Host         string             `json:"host"`
	Port         string             `json:"port"`
	Options      map[string]string  `json:"options"`
	LoggerConfig *MongoLoggerConfig `json:"logger"`
}
