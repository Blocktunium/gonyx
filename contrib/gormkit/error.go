package gormkit

import "fmt"

// NotImplementedErr Error
type NotImplementedErr struct {
}

// Error method - satisfying error interface
func (err *NotImplementedErr) Error() string {
	return fmt.Sprintf("Not Implemented Yet")
}

// NewNotImplementedErr - return a new instance of NotImplementedErr
func NewNotImplementedErr() error {
	return &NotImplementedErr{}
}

// CreateSqlWrapperErr Error
type CreateSqlWrapperErr struct {
	Err error
}

// Error method - satisfying error interface
func (err *CreateSqlWrapperErr) Error() string {
	return fmt.Sprintf("Create a new sql wrapper encounterred an error: %v", err.Err)
}

// NewCreateSqlWrapperErr - return a new instance of CreateSqlWrapperErr
func NewCreateSqlWrapperErr(err error) error {
	return &CreateSqlWrapperErr{Err: err}
}

// NotSupportedDbTypeErr Error
type NotSupportedDbTypeErr struct {
	dbType string
}

// Error method - satisfying error interface
func (err *NotSupportedDbTypeErr) Error() string {
	return fmt.Sprintf("Not Supported Database Dialect: %v", err.dbType)
}

// NewNotSupportedDbTypeErr - return a new instance of NotSupportedDbTypeErr
func NewNotSupportedDbTypeErr(dbType string) error {
	return &NotSupportedDbTypeErr{dbType: dbType}
}

// NotExistServiceNameErr Error
type NotExistServiceNameErr struct {
	serviceName string
}

// Error method - satisfying error interface
func (err *NotExistServiceNameErr) Error() string {
	return fmt.Sprintf("Instance with service name not exist: %v", err.serviceName)
}

// NewNotExistServiceNameErr - return a new instance of NotExistServiceNameErr
func NewNotExistServiceNameErr(serviceName string) error {
	return &NotExistServiceNameErr{serviceName: serviceName}
}

// MigrateErr Error
type MigrateErr struct {
	Err error
}

// Error method - satisfying error interface
func (err *MigrateErr) Error() string {
	return fmt.Sprintf("Migrating tables got error: %v", err.Err)
}

// NewMigrateErr - return a new instance of MigrateErr
func NewMigrateErr(err error) error {
	return &MigrateErr{Err: err}
}

// SelectQueryErr Error
type SelectQueryErr struct {
	query string
	Err   error
}

// Error method - satisfying error interface
func (err *SelectQueryErr) Error() string {
	return fmt.Sprintf("Select query (`%v`) encouters error: %v", err.query, err.Err)
}

// NewSelectQueryErr - return a new instance of SelectQueryErr
func NewSelectQueryErr(q string, err error) error {
	return &SelectQueryErr{query: q, Err: err}
}

// DeleteModelErr Error
type DeleteModelErr struct {
	table string
	data  any
	Err   error
}

// Error method - satisfying error interface
func (err *DeleteModelErr) Error() string {
	return fmt.Sprintf("Deleting a record from (%v) with data: %v -> encouters error: %v", err.table, err.data, err.Err)
}

// NewDeleteModelErr - return a new instance of DeleteModelErr
func NewDeleteModelErr(table string, data any, err error) error {
	return &DeleteModelErr{table: table, data: data, Err: err}
}

// InsertModelErr Error
type InsertModelErr struct {
	table string
	data  any
	Err   error
}

// Error method - satisfying error interface
func (err *InsertModelErr) Error() string {
	return fmt.Sprintf("Inserting a record to (%v) with data: %v -> encouters error: %v", err.table, err.data, err.Err)
}

// NewInsertModelErr - return a new instance of InsertModelErr
func NewInsertModelErr(table string, data any, err error) error {
	return &InsertModelErr{table: table, data: data, Err: err}
}

// UpdateModelErr Error
type UpdateModelErr struct {
	table string
	data  any
	Err   error
}

// Error method - satisfying error interface
func (err *UpdateModelErr) Error() string {
	return fmt.Sprintf("Updating record(s) in (%v) with data: %v -> encouters error: %v", err.table, err.data, err.Err)
}

// NewUpdateModelErr - return a new instance of UpdateModelErr
func NewUpdateModelErr(table string, data any, err error) error {
	return &UpdateModelErr{table: table, data: data, Err: err}
}
