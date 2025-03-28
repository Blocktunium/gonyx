package gormkit

import (
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestSqlWrapper_Initialize(t *testing.T) {
	wrapper := &SqlWrapper[Sqlite]{name: "test"}
	newWrapper, err := NewSqlWrapper[Sqlite]("test", "sqlite")

	if err != nil {
		t.Errorf("Creating Sql Wrapper --> Expected: %v, but got %v", nil, err)
		return
	}

	if !reflect.DeepEqual(wrapper, newWrapper) {
		t.Errorf("Creating Sql Wrapper --> Expected: %v, but got %v", wrapper, newWrapper)
		return
	}
}

func TestSqlWrapper_SqliteConnection(t *testing.T) {
	makeReadyConfigManager()

	newWrapper, err := NewSqlWrapper[Sqlite]("db/server1", "sqlite")
	if err != nil {
		t.Errorf("Creating Sql Wrapper --> Expected: %v, but got %v", nil, err)
		return
	}

	//db, err := gorm.Open(sqlite.Open("file:file.db?mode=memory&cache=shared&_fk=1"), &gorm.Config{})
	//if err != nil {
	//	t.Errorf("Creating Sql Connection Instance --> Expected: %v, but got %v", nil, err)
	//	return
	//}

	db2, err := newWrapper.GetDb()
	if err != nil {
		t.Errorf("Get database instance --> Expected: %v, but got %v", nil, err)
		return
	}

	tx := db2.Exec("SELECT 1")
	if tx.Error != nil {
		t.Errorf("Query on connected database --> Expected no Error: %v, but got %v", nil, err)
		return
	}
}

func TestSqlWrapper_AddLoggerAndTest(t *testing.T) {
	makeReadyConfigManager()

	newWrapper, err := NewSqlWrapper[Sqlite]("db/server1", "sqlite")
	if err != nil {
		t.Errorf("Creating Sql Wrapper --> Expected: %v, but got %v", nil, err)
		return
	}

	db2, err := newWrapper.GetDb()
	if err != nil {
		t.Errorf("Get database instance --> Expected: %v, but got %v", nil, err)
		return
	}

	tx := db2.Exec("SELECT 1")
	if tx.Error != nil {
		t.Errorf("Query on connected database --> Expected no Error: %v, but got %v", nil, err)
		return
	}

}

func TestSqlWrapper_SqliteConnectionConfiguration(t *testing.T) {

	//db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
	//	SkipDefaultTransaction:                   false,
	//	NamingStrategy:                           nil,
	//	FullSaveAssociations:                     false,
	//	Logger:                                   nil,
	//	NowFunc:                                  nil,
	//	DryRun:                                   false,
	//	PrepareStmt:                              false,
	//	DisableAutomaticPing:                     false,
	//	DisableForeignKeyConstraintWhenMigrating: false,
	//	IgnoreRelationshipsWhenMigrating:         false,
	//	DisableNestedTransaction:                 false,
	//	AllowGlobalUpdate:                        false,
	//	QueryFields:                              false,
	//	CreateBatchSize:                          0,
	//	TranslateError:                           false,
	//	ClauseBuilders:                           nil,
	//	ConnPool:                                 nil,
	//	Dialector:                                nil,
	//	Plugins:                                  nil,
	//})
	//
	//db.Logger.LogMode()

}

func TestSqlWrapper_MigrateTest(t *testing.T) {
	makeReadyConfigManager()

	newWrapper, err := NewSqlWrapper[Sqlite]("db/server1", "sqlite")
	if err != nil {
		t.Errorf("Creating Sql Wrapper --> Expected: %v, but got %v", nil, err)
		return
	}

	db2, err := newWrapper.GetDb()
	if err != nil {
		t.Errorf("Get database instance --> Expected: %v, but got %v", nil, err)
		return
	}

	type User struct {
		gorm.Model
		Name string
	}

	err = newWrapper.Migrate(&User{})
	if err != nil {
		t.Errorf("Migrating Tables Error --> Expected: %v, but got %v", nil, err)
		return
	}

	u := User{Name: "test"}

	result := db2.Create(&u) // pass pointer of data to // Create

	if result.Error != nil {
		t.Errorf("Result Error --> Expected: %v, but got %v", nil, result.Error)
		return
	}

	expectedId := uint(1)
	resultId := u.ID
	if resultId != expectedId {
		t.Errorf("Expected To Get First Inserted Id: %v, but got %v", expectedId, resultId)
		return
	}

	expectedRowsAffected := int64(1)
	if result.RowsAffected != expectedRowsAffected {
		t.Errorf("Expected To Get Rows Affected: %v, but got %v", expectedRowsAffected, result.RowsAffected)
		return
	}
}
