package logger

import (
	"github.com/Blocktunium/gonyx/internal/logger"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"time"
)

type LogError logger.Error
type LogObject types.LogObject
type LogLevel types.LogLevel
type LogType types.LogType

// Some Constants - used with LogLevel
const (
	DEBUG   LogLevel = LogLevel(types.DEBUG)
	INFO    LogLevel = LogLevel(types.INFO)
	WARNING LogLevel = LogLevel(types.WARNING)
	ERROR   LogLevel = LogLevel(types.ERROR)
)

var (
	FuncMaintenanceType LogType = LogType(types.NewLogType(types.FuncMaintenanceType.String()))
	DebugType           LogType = LogType(types.NewLogType(types.DebugType.String()))
)

// NewLogObject - enhance method to create and return reference of LogObject
func NewLogObject(level LogLevel, module string, logType LogType, eventTime time.Time, message interface{}, additional interface{}) *LogObject {
	return &LogObject{
		Level:      types.LogLevel(level),
		Module:     module,
		LogType:    types.LogType(logType).String(),
		Time:       eventTime.UTC().UnixNano(),
		Message:    message,
		Additional: additional,
	}
}

// Log - write log object to the channel
func Log(object *LogObject) *LogError {
	l, err := logger.GetManager().GetLogger()
	if err == nil {
		if l.IsInitialized() {
			p := types.LogObject(*object)
			l.Log(&p)
		}
		return nil
	}

	p := LogError(*err)
	return &p
}

// Sync - sync all logs to medium
func Sync() *LogError {
	l, err := logger.GetManager().GetLogger()
	if err == nil {
		if l.IsInitialized() {
			l.Sync()
		}
		return nil
	}

	p := LogError(*err)
	return &p
}

// Close - it closes logger channel
func Close() *LogError {
	l, err := logger.GetManager().GetLogger()
	if err == nil {
		if l.IsInitialized() {
			l.Close()
		}
		return nil
	}

	p := LogError(*err)
	return &p
}
