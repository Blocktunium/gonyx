package extensions

import (
	"fmt"
	gonyxLogger "github.com/Blocktunium/gonyx/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MARK: Variables
var (
	DbLogType      gonyxLogger.LogType
	DbTraceLogType gonyxLogger.LogType
)

func init() {
	// Initialize log types for database operations
	DbLogType = gonyxLogger.NewLogType("DB_OP")
	DbTraceLogType = gonyxLogger.NewLogType("DB_TRACE_OP")
}

// MongoLogger - Mongo DB Logger struct
type MongoLogger struct {
	loggerInstance gonyxLogger.Logger
}

// NewMongoLogger - return instance of MongoLogger which implement Interface
func NewMongoLogger(loggerInstance gonyxLogger.Logger) *MongoLogger {
	return &MongoLogger{
		loggerInstance: loggerInstance,
	}
}

func (m *MongoLogger) Info(level int, message string, keysAndValues ...interface{}) {
	if m.loggerInstance != nil {
		if options.LogLevel(level) == options.LogLevelInfo {
			if logObj := gonyxLogger.NewLogObject(
				gonyxLogger.INFO, "mongodb", DbLogType,
				time.Now().UTC(), message, keysAndValues,
			); logObj != nil {
				m.loggerInstance.Log(logObj)
			}
		} else if options.LogLevel(level) == options.LogLevelDebug {
			if logObj := gonyxLogger.NewLogObject(
				gonyxLogger.DEBUG, "mongodb", DbLogType,
				time.Now().UTC(), message, keysAndValues,
			); logObj != nil {
				m.loggerInstance.Log(logObj)
			}
		}
	}
}

func (m *MongoLogger) Error(err error, message string, keysAndValues ...interface{}) {
	if m.loggerInstance != nil {
		msg := fmt.Sprintf("%s -> err: %v", message, err)
		if logObj := gonyxLogger.NewLogObject(
			gonyxLogger.ERROR, "mongodb", DbLogType,
			time.Now().UTC(), msg, keysAndValues,
		); logObj != nil {
			m.loggerInstance.Log(logObj)
		}
	}
}
