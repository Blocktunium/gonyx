package mongokit

import "fmt"

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

// CreateMongoWrapperErr Error
type CreateMongoWrapperErr struct {
	Err error
}

// Error method - satisfying error interface
func (err *CreateMongoWrapperErr) Error() string {
	return fmt.Sprintf("Create a new mongo wrapper encounterred an error: %v", err.Err)
}

// NewCreateMongoWrapperErr - return a new instance of CreateMongoWrapperErr
func NewCreateMongoWrapperErr(err error) error {
	return &CreateMongoWrapperErr{Err: err}
}

// MongoFindQueryErr Error
type MongoFindQueryErr struct {
	collection string
	filter     any
	Err        error
}

// Error method - satisfying error interface
func (err *MongoFindQueryErr) Error() string {
	return fmt.Sprintf("Find query on (`%s`) with (%v) filter encouters error: %v", err.collection, err.filter, err.Err)
}

// NewMongoFindQueryErr - return a new instance of MongoFindQueryErr
func NewMongoFindQueryErr(collection string, filter any, err error) error {
	return &MongoFindQueryErr{collection: collection, filter: filter, Err: err}
}

// MongoDeleteErr Error
type MongoDeleteErr struct {
	collection string
	filter     any
	Err        error
}

// Error method - satisfying error interface
func (err *MongoDeleteErr) Error() string {
	return fmt.Sprintf("Delete query on (`%s`) with (%v) filter encouters error: %v", err.collection, err.filter, err.Err)
}

// NewMongoDeleteErr - return a new instance of MongoDeleteErr
func NewMongoDeleteErr(collection string, filter any, err error) error {
	return &MongoDeleteErr{collection: collection, filter: filter, Err: err}
}
