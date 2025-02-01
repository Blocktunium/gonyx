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
