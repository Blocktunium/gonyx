package extensions

import (
	"context"
	"errors"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

// MARK: Variables
var (
	DbLogType      = types.NewLogType("DB_OP")
	DbTraceLogType = types.NewLogType("DB_TRACE_OP")
)

// DbLogger - DB Logger struct
type DbLogger struct {
	logger.Config
	loggerInstance types.Logger
}

// NewDbLogger - return instance of DbLogger which implement Interface
func NewDbLogger(config map[string]interface{}, loggerInstance types.Logger) logger.Interface {
	logLevel := logger.Silent
	switch strings.ToLower(config["log_level"].(string)) {
	case "error":
		logLevel = logger.Error
		break
	case "warn":
		logLevel = logger.Warn
		break
	case "info":
		logLevel = logger.Info
		break
	}

	l := &DbLogger{Config: logger.Config{
		SlowThreshold:             time.Duration(int(config["slow_threshold"].(float64))) * time.Millisecond,
		Colorful:                  false,
		IgnoreRecordNotFoundError: config["ignore_record_not_found_error"].(bool),
		ParameterizedQueries:      config["parameterized_queries"].(bool),
		LogLevel:                  logLevel,
	}, loggerInstance: loggerInstance}

	return l
}

// LogMode - set log mode
func (l *DbLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info - print info
func (l DbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info && l.loggerInstance != nil {
		newMsg := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		l.loggerInstance.Log(types.NewLogObject(
			types.INFO, "db", DbLogType,
			time.Now().UTC(), newMsg, nil,
		))
	}
}

// Warn - print warn messages
func (l DbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn && l.loggerInstance != nil {
		newMsg := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		l.loggerInstance.Log(types.NewLogObject(
			types.WARNING, "db", DbLogType,
			time.Now().UTC(), newMsg, nil,
		))
	}
}

// Error - print error messages
func (l DbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error && l.loggerInstance != nil {
		newMsg := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
		l.loggerInstance.Log(types.NewLogObject(
			types.ERROR, "db", DbLogType,
			time.Now().UTC(), newMsg, nil,
		))
	}
}

// Trace - print sql message
func (l DbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	if l.loggerInstance != nil {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= logger.Error &&
			(!errors.Is(err, logger.ErrRecordNotFound) ||
				!l.IgnoreRecordNotFoundError):
			sql, rows := fc()
			msgLiteral := "%s %s\n[%.3fms] [rows:%v] %s"
			if rows == -1 {
				msg := fmt.Sprintf(msgLiteral, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
				l.loggerInstance.Log(types.NewLogObject(
					types.ERROR, "db", DbTraceLogType,
					time.Now().UTC(), msg, []interface{}{err, sql},
				))
			} else {
				msg := fmt.Sprintf(msgLiteral, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
				l.loggerInstance.Log(types.NewLogObject(
					types.ERROR, "db", DbTraceLogType,
					time.Now().UTC(), msg, []interface{}{err, sql, rows},
				))
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)

			msgLiteral := "%s %s\n[%.3fms] [rows:%v] %s"
			if rows == -1 {
				msg := fmt.Sprintf(msgLiteral, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
				l.loggerInstance.Log(types.NewLogObject(
					types.WARNING, "db", DbTraceLogType,
					time.Now().UTC(), msg, []interface{}{slowLog, sql},
				))
			} else {
				msg := fmt.Sprintf(msgLiteral, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
				l.loggerInstance.Log(types.NewLogObject(
					types.WARNING, "db", DbTraceLogType,
					time.Now().UTC(), msg, []interface{}{slowLog, sql, rows},
				))
			}
		case l.LogLevel == logger.Info:
			sql, rows := fc()

			msgLiteral := "%s\n[%.3fms] [rows:%v] %s"
			if rows == -1 {
				msg := fmt.Sprintf(msgLiteral, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
				l.loggerInstance.Log(types.NewLogObject(
					types.INFO, "db", DbTraceLogType,
					time.Now().UTC(), msg, []interface{}{sql},
				))
			} else {
				msg := fmt.Sprintf(msgLiteral, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
				l.loggerInstance.Log(types.NewLogObject(
					types.INFO, "db", DbTraceLogType,
					time.Now().UTC(), msg, []interface{}{sql, rows},
				))
			}
		}
	}
}

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
