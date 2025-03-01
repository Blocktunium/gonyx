package logger

// Imports needed list
import (
	"errors"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"github.com/Blocktunium/gonyx/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Mark: ZapWrapper

// ZapWrapper structure - implements Logger interface
type ZapWrapper struct {
	name            string
	serviceName     string
	logger          *zap.Logger
	ch              chan types.LogObject
	initialized     bool
	wg              sync.WaitGroup
	operationType   string
	supportedOutput []string
}

// Constructor - It initializes the logger configuration params
func (l *ZapWrapper) Constructor(name string) error {
	l.wg.Add(1)
	defer l.wg.Done()

	l.name = name
	l.serviceName = config.GetManager().GetName()
	l.operationType = config.GetManager().GetOperationType()
	l.supportedOutput = []string{"console", "file"}
	l.initialized = false

	channelSize, err := config.GetManager().Get(l.name, "channel_size")
	if err != nil {
		return err
	}

	options, err := config.GetManager().Get(l.name, "options")
	if err != nil {
		return err
	}

	var optionArray []string
	for _, v := range options.([]interface{}) {
		optionArray = append(optionArray, v.(string))
	}

	outputs, err := config.GetManager().Get(l.name, "outputs")
	if err != nil {
		return err
	}

	var outputArray []string
	for _, v := range outputs.([]interface{}) {
		outputArray = append(outputArray, v.(string))
	}

	l.ch = make(chan types.LogObject, int(channelSize.(float64)))

	if l.operationType == "prod" {
		productionEncoderConfig := zap.NewProductionEncoderConfig()
		productionEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		productionEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		var cores []zapcore.Core
		for _, outputItem := range outputArray {
			if utils.ArrayContains(&l.supportedOutput, outputItem) {
				if outputItem == "console" {
					level := zapcore.DebugLevel

					key := fmt.Sprintf("%s.level", outputItem)
					levelStr, err := config.GetManager().Get(l.name, key)
					if err == nil {
						level, err = zapcore.ParseLevel(levelStr.(string))
						if err != nil {
							continue
						}
					}

					consoleEncoder := zapcore.NewConsoleEncoder(productionEncoderConfig)
					c := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
					cores = append(cores, c)
				} else if outputItem == "file" {
					level := zapcore.DebugLevel

					key := fmt.Sprintf("%s.level", outputItem)
					levelStr, err := config.GetManager().Get(l.name, key)
					if err == nil {
						level, err = zapcore.ParseLevel(levelStr.(string))
						if err != nil {
							continue
						}
					}

					// Read the root path of logs
					path := "logs"
					key = fmt.Sprintf("%s.path", outputItem)
					pathStr, err := config.GetManager().Get(l.name, key)
					if err == nil {
						if strings.TrimSpace(pathStr.(string)) != "" {
							path = strings.TrimSpace(pathStr.(string))
						}
					}

					// Check the directory existed, If not create all the nested directories
					if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
						err1 := os.MkdirAll(path, os.ModePerm)
						if err1 != nil {
							continue
						}
					}

					expectLogPath := filepath.Join(path, fmt.Sprintf("%s.log", config.GetManager().GetName()))
					logFile, err := os.OpenFile(expectLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
					if err != nil {
						continue
					}
					writer := zapcore.AddSync(logFile)
					fileEncoder := zapcore.NewJSONEncoder(productionEncoderConfig)

					c := zapcore.NewCore(fileEncoder, writer, level)
					cores = append(cores, c)
				}
			}
		}

		core := zapcore.NewTee(
			cores...,
		)

		if utils.ArrayContains(&optionArray, "stackTrace") && utils.ArrayContains(&optionArray, "caller") {
			l.logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCaller())
		} else if utils.ArrayContains(&optionArray, "stackTrace") {
			l.logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
		} else if utils.ArrayContains(&optionArray, "caller") {
			l.logger = zap.New(core, zap.AddCaller())
		} else {
			l.logger = zap.New(core)
		}
	} else { // otherwise we create a development version (`dev`, `test`, ...)
		developmentEncoderConfig := zap.NewDevelopmentEncoderConfig()
		developmentEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		developmentEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		var cores []zapcore.Core
		for _, outputItem := range outputArray {
			if utils.ArrayContains(&l.supportedOutput, outputItem) {
				if outputItem == "console" {
					level := zapcore.DebugLevel

					key := fmt.Sprintf("%s.level", outputItem)
					levelStr, err := config.GetManager().Get(l.name, key)
					if err == nil {
						level, err = zapcore.ParseLevel(levelStr.(string))
						if err != nil {
							continue
						}
					}

					consoleEncoder := zapcore.NewConsoleEncoder(developmentEncoderConfig)
					c := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)

					cores = append(cores, c)
				} else if outputItem == "file" {
					level := zapcore.DebugLevel

					key := fmt.Sprintf("%s.level", outputItem)
					levelStr, err := config.GetManager().Get(l.name, key)
					if err == nil {
						level, err = zapcore.ParseLevel(levelStr.(string))
						if err != nil {
							continue
						}
					}

					// Read the root path of logs
					path := "logs"
					key = fmt.Sprintf("%s.path", outputItem)
					pathStr, err := config.GetManager().Get(l.name, key)
					if err == nil {
						if strings.TrimSpace(pathStr.(string)) != "" {
							path = strings.TrimSpace(pathStr.(string))
						}
					}

					// Check the directory existed, If not create all the nested directories
					if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
						err1 := os.MkdirAll(path, os.ModePerm)
						if err1 != nil {
							continue
						}
					}

					expectLogPath := filepath.Join(path, fmt.Sprintf("%s.log", config.GetManager().GetName()))
					logFile, err := os.OpenFile(expectLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
					if err != nil {
						continue
					}
					writer := zapcore.AddSync(logFile)
					fileEncoder := zapcore.NewJSONEncoder(developmentEncoderConfig)

					c := zapcore.NewCore(fileEncoder, writer, level)
					cores = append(cores, c)
				}
			}
		}

		core := zapcore.NewTee(
			cores...,
		)

		if utils.ArrayContains(&optionArray, "stackTrace") && utils.ArrayContains(&optionArray, "caller") {
			l.logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCaller())
		} else if utils.ArrayContains(&optionArray, "stackTrace") {
			l.logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
		} else if utils.ArrayContains(&optionArray, "caller") {
			l.logger = zap.New(core, zap.AddCaller())
		} else {
			l.logger = zap.New(core)
		}
	}

	go l.runner()
	l.initialized = true

	return nil
}

// Close - it closes logger channel
func (l *ZapWrapper) Close() {
	l.wg.Wait()

	_ = l.logger.Sync()
	defer close(l.ch)
}

// Log - write log object to the channel
func (l *ZapWrapper) Log(obj *types.LogObject) {
	l.wg.Wait()

	go func(obj *types.LogObject) {
		l.ch <- *obj
	}(obj)
}

// IsInitialized - that returns boolean value whether it's initialized
func (l *ZapWrapper) IsInitialized() bool {
	return l.initialized
}

// Instance - returns exact logger instance
func (l *ZapWrapper) Instance() *zap.Logger {
	l.wg.Wait()
	return l.logger
}

// Sync - call the sync method of the project
func (l *ZapWrapper) Sync() {
	l.wg.Wait()
	l.logger.Sync()
}

// runner - the goroutine that reads from channel and process it
func (l *ZapWrapper) runner() {
	l.wg.Wait()
	for c := range l.ch {
		if c.Level.IsLogLevel() {
			f := []zapcore.Field{
				zap.Any("service", l.serviceName),
				zap.Any("module", c.Module),
				zap.Any("log_type", c.LogType),
				zap.Any("time", c.Time),
				zap.Any("additional", c.Additional),
			}
			switch c.Level {
			case types.DEBUG:
				l.logger.Debug(fmt.Sprintf("%v", c.Message), f...)
				break
			case types.INFO:
				l.logger.Info(fmt.Sprintf("%v", c.Message), f...)
				break
			case types.WARNING:
				l.logger.Warn(fmt.Sprintf("%v", c.Message), f...)
				break
			case types.ERROR:
				l.logger.Error(fmt.Sprintf("%v", c.Message), f...)
				break
			}
		}
	}
}
