package rediskit

import "fmt"

// Error object
type Error struct {
	Err error
}

// Error method - satisfying error interface
func (err *Error) Error() string {
	return fmt.Sprintf("Redis Error | %v", err.Err.Error())
}

// NewError - return a new instance of Error
func NewError(err error) error {
	return &Error{Err: err}
}

// ReadError struct
type ReadError struct {
	Key string
	Err error
}

// Error method - satisfying error interface
func (err *ReadError) Error() string {
	return fmt.Sprintf("Redis Read Error: key = %s | %v", err.Key, err.Err)
}

// NewReadError - return a new instance of ReadError
func NewReadError(key string, err error) error {
	return &ReadError{
		Key: key,
		Err: err,
	}
}

// WriteError struct
type WriteError struct {
	Key   string
	Value any
	Err   error
}

// Error method - satisfying error interface
func (err *WriteError) Error() string {
	return fmt.Sprintf("Redis Write Error: key = %s / value = %v | %v", err.Key, err.Value, err.Err)
}

// NewWriteError - return a new instance of WriteError
func NewWriteError(key string, value any, err error) error {
	return &WriteError{
		Key:   key,
		Value: value,
		Err:   err,
	}
}

// PingError struct
type PingError struct {
	Err error
}

// Error method - satisfying error interface
func (err *PingError) Error() string {
	return fmt.Sprintf("Redis Ping Error | %v", err.Err)
}

// NewPingError - return a new instance of PingError
func NewPingError(err error) error {
	return &PingError{Err: err}
}
