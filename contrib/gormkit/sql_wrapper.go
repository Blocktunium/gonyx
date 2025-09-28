package gormkit

import (
	"encoding/json"
	"fmt"
	"github.com/Blocktunium/gonyx-gormkit/extensions"
	"github.com/Blocktunium/gonyx/pkg/config"
	"github.com/Blocktunium/gonyx/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

// Mark: Definitions

// SqlWrapper struct
type SqlWrapper[T SqlConfigurable] struct {
	name             string
	config           T
	databaseInstance *gorm.DB
	logger           logger.Logger // Gonyx logger interface
}

// init - SqlWrapper Constructor - It initializes the wrapper
func (s *SqlWrapper[T]) init(name string) error {
	s.name = name

	// reading config
	nameParts := strings.Split(s.name, "/")

	if reflect.ValueOf(s.config).Type() == reflect.TypeOf(Sqlite{}) {
		filenameKey := fmt.Sprintf("%s.%s", nameParts[1], "db")
		filenameStr, err := config.Get(nameParts[0], filenameKey)
		if err != nil {
			return err
		}

		optionsKey := fmt.Sprintf("%s.%s", nameParts[1], "options")
		optionsObj, err := config.Get(nameParts[0], optionsKey)
		if err != nil {
			return err
		}

		optionsMap := make(map[string]string, len(optionsObj.(map[string]interface{})))
		for key, item := range optionsObj.(map[string]interface{}) {
			optionsMap[key] = item.(string)
		}

		var internalConfig *Config

		internalConfigKey := fmt.Sprintf("%s.%s", nameParts[1], "config")
		internalConfigObj, err := config.Get(nameParts[0], internalConfigKey)
		if err == nil {
			// first marshal
			configData, err := json.Marshal(internalConfigObj)
			if err == nil {
				_ = json.Unmarshal(configData, &internalConfig)
			}
		}

		var internalLogger *LoggerConfig

		internalLoggerKey := fmt.Sprintf("%s.%s", nameParts[1], "logger")
		internalLoggerObj, err := config.Get(nameParts[0], internalLoggerKey)
		if err == nil {
			// first marshal
			configData, err := json.Marshal(internalLoggerObj)
			if err == nil {
				_ = json.Unmarshal(configData, &internalLogger)
			}
		}

		s.config = reflect.ValueOf(Sqlite{
			FileName:     filenameStr.(string),
			Options:      optionsMap,
			Config:       internalConfig,
			LoggerConfig: internalLogger,
		}).Interface().(T)
	} else if reflect.ValueOf(s.config).Type() == reflect.TypeOf(Mysql{}) {
		dbNameKey := fmt.Sprintf("%s.%s", nameParts[1], "db")
		dbNameStr, err := config.Get(nameParts[0], dbNameKey)
		if err != nil {
			return err
		}

		hostKey := fmt.Sprintf("%s.%s", nameParts[1], "host")
		hostStr, err := config.Get(nameParts[0], hostKey)
		if err != nil {
			return err
		}

		portKey := fmt.Sprintf("%s.%s", nameParts[1], "port")
		portStr, err := config.Get(nameParts[0], portKey)
		if err != nil {
			return err
		}

		protocolKey := fmt.Sprintf("%s.%s", nameParts[1], "protocol")
		protocolStr, err := config.Get(nameParts[0], protocolKey)
		if err != nil {
			return err
		}

		usernameKey := fmt.Sprintf("%s.%s", nameParts[1], "username")
		usernameStr, err := config.Get(nameParts[0], usernameKey)
		if err != nil {
			return err
		}

		passwordKey := fmt.Sprintf("%s.%s", nameParts[1], "password")
		passwordStr, err := config.Get(nameParts[0], passwordKey)
		if err != nil {
			return err
		}

		optionsKey := fmt.Sprintf("%s.%s", nameParts[1], "options")
		optionsObj, err := config.Get(nameParts[0], optionsKey)
		if err != nil {
			return err
		}

		optionsMap := make(map[string]string, len(optionsObj.(map[string]interface{})))
		for key, item := range optionsObj.(map[string]interface{}) {
			optionsMap[key] = item.(string)
		}

		var internalConfig *Config

		internalConfigKey := fmt.Sprintf("%s.%s", nameParts[1], "config")
		internalConfigObj, err := config.Get(nameParts[0], internalConfigKey)
		if err == nil {
			// first marshal
			configData, err := json.Marshal(internalConfigObj)
			if err == nil {
				_ = json.Unmarshal(configData, &internalConfig)
			}
		}

		var internalLogger *LoggerConfig

		internalLoggerKey := fmt.Sprintf("%s.%s", nameParts[1], "logger")
		internalLoggerObj, err := config.Get(nameParts[0], internalLoggerKey)
		if err == nil {
			// first marshal
			configData, err := json.Marshal(internalLoggerObj)
			if err == nil {
				_ = json.Unmarshal(configData, &internalLogger)
			}
		}

		var specificConfig *MysqlSpecificConfig

		specificConfigKey := fmt.Sprintf("%s.%s", nameParts[1], "specific_config")
		specificConfigObj, err := config.Get(nameParts[0], specificConfigKey)
		if err == nil {
			// first marshal
			configData, err := json.Marshal(specificConfigObj)
			if err == nil {
				_ = json.Unmarshal(configData, &specificConfig)
			}
		}

		s.config = reflect.ValueOf(Mysql{
			DatabaseName:   dbNameStr.(string),
			Username:       usernameStr.(string),
			Password:       passwordStr.(string),
			Host:           hostStr.(string),
			Port:           portStr.(string),
			Protocol:       protocolStr.(string),
			Options:        optionsMap,
			Config:         internalConfig,
			LoggerConfig:   internalLogger,
			SpecificConfig: specificConfig,
		}).Interface().(T)
	} else if reflect.ValueOf(s.config).Type() == reflect.TypeOf(Postgresql{}) {
		dbNameKey := fmt.Sprintf("%s.%s", nameParts[1], "db")
		dbNameStr, err := config.Get(nameParts[0], dbNameKey)
		if err != nil {
			return err
		}

		hostKey := fmt.Sprintf("%s.%s", nameParts[1], "host")
		hostStr, err := config.Get(nameParts[0], hostKey)
		if err != nil {
			return err
		}

		portKey := fmt.Sprintf("%s.%s", nameParts[1], "port")
		portStr, err := config.Get(nameParts[0], portKey)
		if err != nil {
			return err
		}

		usernameKey := fmt.Sprintf("%s.%s", nameParts[1], "username")
		usernameStr, err := config.Get(nameParts[0], usernameKey)
		if err != nil {
			return err
		}

		passwordKey := fmt.Sprintf("%s.%s", nameParts[1], "password")
		passwordStr, err := config.Get(nameParts[0], passwordKey)
		if err != nil {
			return err
		}

		optionsKey := fmt.Sprintf("%s.%s", nameParts[1], "options")
		optionsObj, err := config.Get(nameParts[0], optionsKey)
		if err != nil {
			return err
		}

		optionsMap := make(map[string]string, len(optionsObj.(map[string]interface{})))
		for key, item := range optionsObj.(map[string]interface{}) {
			optionsMap[key] = item.(string)
		}

		var internalConfig *Config

		internalConfigKey := fmt.Sprintf("%s.%s", nameParts[1], "config")
		internalConfigObj, err := config.Get(nameParts[0], internalConfigKey)
		if err == nil {
			// first marshal
			configData, err := json.Marshal(internalConfigObj)
			if err == nil {
				_ = json.Unmarshal(configData, &internalConfig)
			}
		}

		var internalLogger *LoggerConfig

		internalLoggerKey := fmt.Sprintf("%s.%s", nameParts[1], "logger")
		internalLoggerObj, err := config.Get(nameParts[0], internalLoggerKey)
		if err == nil {
			// first marshal
			configData, err := json.Marshal(internalLoggerObj)
			if err == nil {
				_ = json.Unmarshal(configData, &internalLogger)
			}
		}

		var specificConfig *PostgresqlSpecificConfig

		specificConfigKey := fmt.Sprintf("%s.%s", nameParts[1], "specific_config")
		specificConfigObj, err := config.Get(nameParts[0], specificConfigKey)
		if err == nil {
			// first marshal
			configData, err := json.Marshal(specificConfigObj)
			if err == nil {
				_ = json.Unmarshal(configData, &specificConfig)
			}
		}

		s.config = reflect.ValueOf(Postgresql{
			DatabaseName:   dbNameStr.(string),
			Username:       usernameStr.(string),
			Password:       passwordStr.(string),
			Host:           hostStr.(string),
			Port:           portStr.(string),
			Options:        optionsMap,
			Config:         internalConfig,
			LoggerConfig:   internalLogger,
			SpecificConfig: specificConfig,
		}).Interface().(T)
	}

	return nil
}

// MARK: Public functions

// RegisterLogger - register logger instance
func (s *SqlWrapper[T]) RegisterLogger(l logger.Logger) {
	s.logger = l
}

// GetDb - return associated internal Db
func (s *SqlWrapper[T]) GetDb() (*gorm.DB, error) {
	if s.databaseInstance == nil {
		if reflect.ValueOf(s.config).Type() == reflect.TypeOf(Sqlite{}) {
			optionsQSArr := make([]string, 0)
			config := reflect.ValueOf(s.config).Interface().(Sqlite)
			for key, val := range config.Options {
				optionsQSArr = append(optionsQSArr, fmt.Sprintf("%s=%s", key, val))
			}
			optionsQS := strings.Join(optionsQSArr, "&")

			dsn := fmt.Sprintf("file:%s?%s", config.FileName, optionsQS)
			internalConfig := &gorm.Config{}
			if config.Config != nil {
				internalConfig.DisableAutomaticPing = config.Config.DisableAutomaticPing
				internalConfig.DisableForeignKeyConstraintWhenMigrating = config.Config.DisableForeignKeyConstraintWhenMigrating
				internalConfig.DisableNestedTransaction = config.Config.DisableNestedTransaction
				internalConfig.DryRun = config.Config.DryRun
				internalConfig.PrepareStmt = config.Config.PrepareStmt
				internalConfig.SkipDefaultTransaction = config.Config.SkipDefaultTransaction
				internalConfig.IgnoreRelationshipsWhenMigrating = config.Config.IgnoreRelationshipsWhenMigrating
			}

			if config.LoggerConfig != nil {
				var inInterface map[string]interface{}
				inrec, _ := json.Marshal(config.LoggerConfig)
				json.Unmarshal(inrec, &inInterface)
				internalConfig.Logger = extensions.NewDbLogger(inInterface, s.logger)
			}

			db, err := gorm.Open(sqlite.Open(dsn), internalConfig)
			if err != nil {
				return nil, err
			}
			s.databaseInstance = db
		} else if reflect.ValueOf(s.config).Type() == reflect.TypeOf(Mysql{}) {
			optionsQSArr := make([]string, 0)
			config := reflect.ValueOf(s.config).Interface().(Mysql)
			for key, val := range config.Options {
				optionsQSArr = append(optionsQSArr, fmt.Sprintf("%s=%s", key, val))
			}
			optionsQS := strings.Join(optionsQSArr, "&")

			dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?%s", config.Username,
				config.Password, config.Protocol, config.Host, config.Port,
				config.DatabaseName, optionsQS)
			internalConfig := &gorm.Config{}
			if config.Config != nil {
				internalConfig.DisableAutomaticPing = config.Config.DisableAutomaticPing
				internalConfig.DisableForeignKeyConstraintWhenMigrating = config.Config.DisableForeignKeyConstraintWhenMigrating
				internalConfig.DisableNestedTransaction = config.Config.DisableNestedTransaction
				internalConfig.DryRun = config.Config.DryRun
				internalConfig.PrepareStmt = config.Config.PrepareStmt
				internalConfig.SkipDefaultTransaction = config.Config.SkipDefaultTransaction
				internalConfig.IgnoreRelationshipsWhenMigrating = config.Config.IgnoreRelationshipsWhenMigrating
			}

			if config.LoggerConfig != nil {
				var inInterface map[string]interface{}
				inrec, _ := json.Marshal(config.LoggerConfig)
				json.Unmarshal(inrec, &inInterface)

				internalConfig.Logger = extensions.NewDbLogger(inInterface, s.logger)
			}

			if config.SpecificConfig == nil {
				db, err := gorm.Open(mysql.Open(dsn), internalConfig)
				if err != nil {
					return nil, err
				}
				s.databaseInstance = db
			} else {
				db, err := gorm.Open(mysql.New(mysql.Config{
					DSN:                           dsn,
					SkipInitializeWithVersion:     config.SpecificConfig.SkipInitializeWithVersion,
					DefaultStringSize:             config.SpecificConfig.DefaultStringSize,
					DefaultDatetimePrecision:      &config.SpecificConfig.DefaultDatetimePrecision,
					DisableWithReturning:          config.SpecificConfig.DisableWithReturning,
					DisableDatetimePrecision:      config.SpecificConfig.DisableDatetimePrecision,
					DontSupportRenameIndex:        !config.SpecificConfig.SupportRenameIndex,
					DontSupportRenameColumn:       !config.SpecificConfig.SupportRenameColumn,
					DontSupportForShareClause:     !config.SpecificConfig.SupportForShareClause,
					DontSupportNullAsDefaultValue: !config.SpecificConfig.SupportNullAsDefaultValue,
					DontSupportRenameColumnUnique: !config.SpecificConfig.SupportRenameColumnUnique,
				}), internalConfig)
				if err != nil {
					return nil, err
				}
				s.databaseInstance = db
			}

		} else if reflect.ValueOf(s.config).Type() == reflect.TypeOf(Postgresql{}) {
			optionsQSArr := make([]string, 0)
			config := reflect.ValueOf(s.config).Interface().(Postgresql)
			for key, val := range config.Options {
				optionsQSArr = append(optionsQSArr, fmt.Sprintf("%s=%s", key, val))
			}
			optionsQS := strings.Join(optionsQSArr, " ")

			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s",
				config.Host, config.Username, config.Password, config.DatabaseName,
				config.Port, optionsQS,
			)
			internalConfig := &gorm.Config{}
			if config.Config != nil {
				internalConfig.DisableAutomaticPing = config.Config.DisableAutomaticPing
				internalConfig.DisableForeignKeyConstraintWhenMigrating = config.Config.DisableForeignKeyConstraintWhenMigrating
				internalConfig.DisableNestedTransaction = config.Config.DisableNestedTransaction
				internalConfig.DryRun = config.Config.DryRun
				internalConfig.PrepareStmt = config.Config.PrepareStmt
				internalConfig.SkipDefaultTransaction = config.Config.SkipDefaultTransaction
				internalConfig.IgnoreRelationshipsWhenMigrating = config.Config.IgnoreRelationshipsWhenMigrating
			}

			if config.LoggerConfig != nil {
				var inInterface map[string]interface{}
				inrec, _ := json.Marshal(config.LoggerConfig)
				json.Unmarshal(inrec, &inInterface)

				internalConfig.Logger = extensions.NewDbLogger(inInterface, s.logger)
			}

			if config.SpecificConfig == nil {
				db, err := gorm.Open(postgres.Open(dsn), internalConfig)
				if err != nil {
					return nil, err
				}
				s.databaseInstance = db
			} else {
				db, err := gorm.Open(postgres.New(postgres.Config{
					DSN:                  dsn,
					PreferSimpleProtocol: config.SpecificConfig.PreferSimpleProtocol,
					WithoutReturning:     config.SpecificConfig.WithoutReturning,
				}), internalConfig)
				if err != nil {
					return nil, err
				}

				x, err := db.DB()
				if err == nil {
					x.SetMaxIdleConns(int(config.SpecificConfig.MaxIdleConnCount))
					x.SetMaxOpenConns(int(config.SpecificConfig.MaxOpenConnCount))
					x.SetConnMaxIdleTime(time.Minute * time.Duration(config.SpecificConfig.ConnMaxLifetime))
				}
				s.databaseInstance = db

			}
		}
	}
	return s.databaseInstance, nil
}

// Migrate - migrate models to the database
func (s *SqlWrapper[T]) Migrate(models ...interface{}) error {
	db, err := s.GetDb()
	if err != nil {
		return err
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return NewMigrateErr(err)
	}
	return nil
}

// AttachMigrationFunc -  attach migration function to be called by end user
func (s *SqlWrapper[T]) AttachMigrationFunc(f func(migrator gorm.Migrator) error) error {
	err := f(s.databaseInstance.Migrator())
	if err != nil {
		return NewMigrateErr(err)
	}
	return nil
}

// NewSqlWrapper - create a new instance of SqlWrapper and returns it
func NewSqlWrapper[T SqlConfigurable](name string, dbType string) (*SqlWrapper[T], error) {
	if strings.ToLower(dbType) == "sqlite" ||
		strings.ToLower(dbType) == "mysql" ||
		strings.ToLower(dbType) == "postgresql" {
		wrapper := &SqlWrapper[T]{}
		err := wrapper.init(name)
		if err != nil {
			return nil, NewCreateSqlWrapperErr(err)
		}

		return wrapper, nil
	}

	return nil, NewNotSupportedDbTypeErr(dbType)
}
