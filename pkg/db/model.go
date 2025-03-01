package db

import (
	"github.com/Blocktunium/gonyx/internal/db"
	"github.com/Blocktunium/gonyx/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// GetDb - Get *gorm.DB instance from the underlying interfaces
func GetDb(instanceName string) (*gorm.DB, error) {
	return db.GetManager().GetDb(instanceName)
}

// Migrate - migrate models on specific database
func Migrate(instanceName string, models ...interface{}) error {
	return db.GetManager().Migrate(instanceName, models...)
}

// AttachMigrationFunc -  attach migration function to be called by end user
func AttachMigrationFunc(instanceName string, f func(migrator gorm.Migrator) error) error {
	return db.GetManager().AttachMigrationFunc(instanceName, f)
}

// GetMongoDb - Get *mongo.Client instance from the underlying interfaces
func GetMongoDb(instanceName string) (*mongo.Database, error) {
	return db.GetManager().GetMongoDb(instanceName)
}

// SetupManager - Setup manager
func SetupManager() {
	l, _ := logger.GetManager().GetLogger()
	db.GetManager().RegisterLogger(l)
}

// MARK: TYPE ALIASES

type InsertModelErr = db.InsertModelErr
type UpdateModelErr = db.UpdateModelErr
type SelectQueryErr = db.SelectQueryErr
type DeleteModelErr = db.DeleteModelErr
type MongoFindQueryErr = db.MongoFindQueryErr
type MongoDeleteErr = db.MongoDeleteErr

// NewInsertModelErr - return a new instance of InsertModelErr
func NewInsertModelErr(table string, data any, err error) error {
	return db.NewInsertModelErr(table, data, err)
}

// NewUpdateModelErr - return a new instance of UpdateModelErr
func NewUpdateModelErr(table string, data any, err error) error {
	return db.NewUpdateModelErr(table, data, err)
}

// NewDeleteModelErr - return a new instance of DeleteModelErr
func NewDeleteModelErr(table string, data any, err error) error {
	return db.NewDeleteModelErr(table, data, err)
}

// NewSelectQueryErr - return a new instance of SelectQueryErr
func NewSelectQueryErr(query string, err error) error {
	return db.NewSelectQueryErr(query, err)
}

// NewMongoFindQueryErr - return a new instance of MongoFindQueryErr
func NewMongoFindQueryErr(collection string, filter any, err error) error {
	return db.NewMongoFindQueryErr(collection, filter, err)
}

// NewMongoDeleteErr - return a new instance of MongoDeleteErr
func NewMongoDeleteErr(collection string, filter any, err error) error {
	return db.NewMongoDeleteErr(collection, filter, err)
}
