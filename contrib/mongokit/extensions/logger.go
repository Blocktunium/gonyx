package extensions

import (
	"fmt"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MARK: Variables
var (
	DbLogType      = types.NewLogType("DB_OP")
	DbTraceLogType = types.NewLogType("DB_TRACE_OP")
)

// MongoLogger - Mongo DB Logger struct
type MongoLogger struct {
	loggerInstance types.Logger
}

// NewMongoLogger - return instance of MongoLogger which implement Interface
func NewMongoLogger(loggerInstance types.Logger) *MongoLogger {
	return &MongoLogger{
		loggerInstance: loggerInstance,
	}
}

func (m *MongoLogger) Info(level int, message string, keysAndValues ...interface{}) {
	if m.loggerInstance != nil {
		if options.LogLevel(level) == options.LogLevelInfo {
			m.loggerInstance.Log(types.NewLogObject(
				types.INFO, "mongodb", DbLogType,
				time.Now().UTC(), message, keysAndValues,
			))
		} else if options.LogLevel(level) == options.LogLevelDebug {
			m.loggerInstance.Log(types.NewLogObject(
				types.DEBUG, "mongodb", DbLogType,
				time.Now().UTC(), message, keysAndValues,
			))
		}
	}
}

func (m *MongoLogger) Error(err error, message string, keysAndValues ...interface{}) {
	if m.loggerInstance != nil {
		msg := fmt.Sprintf("%s -> err: %v", message, err)
		m.loggerInstance.Log(types.NewLogObject(
			types.ERROR, "mongodb", DbLogType,
			time.Now().UTC(), msg, keysAndValues,
		))
	}
}
