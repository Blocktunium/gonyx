package helpers

import (
	"github.com/Blocktunium/gonyx/internal/db"
	"gorm.io/gorm"
)

func GetSqlDbInstance(instanceName string) (*gorm.DB, error) {
	return db.GetManager().GetDb(instanceName)
}
