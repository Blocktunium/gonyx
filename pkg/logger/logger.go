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

// Logger interface - using pkg types instead of internal types
type Logger interface {
	Constructor(name string) error
	Close()
	Log(obj *LogObject)
	IsInitialized() bool
	Sync()
}

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

// NewLogType - Create a new log type
func NewLogType(name string) LogType {
	return LogType(types.NewLogType(name))
}

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

// loggerAdapter adapts internal logger to pkg Logger interface
type loggerAdapter struct {
	internal types.Logger
}

func (la *loggerAdapter) Constructor(name string) error {
	return la.internal.Constructor(name)
}

func (la *loggerAdapter) Close() {
	la.internal.Close()
}

func (la *loggerAdapter) Log(obj *LogObject) {
	// Convert pkg LogObject to internal types LogObject
	internalObj := &types.LogObject{
		Level:      types.LogLevel(obj.Level),
		Module:     obj.Module,
		LogType:    obj.LogType,
		Time:       obj.Time,
		Additional: obj.Additional,
		Message:    obj.Message,
	}
	la.internal.Log(internalObj)
}

func (la *loggerAdapter) IsInitialized() bool {
	return la.internal.IsInitialized()
}

func (la *loggerAdapter) Sync() {
	la.internal.Sync()
}

// GetLogger - returns the logger interface
func GetLogger() (Logger, *LogError) {
	l, err := logger.GetManager().GetLogger()
	if err == nil {
		return &loggerAdapter{internal: l}, nil
	}

	p := LogError(*err)
	return nil, &p
}
