package logger

// Imports needed list
import (
	"encoding/json"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"github.com/Blocktunium/gonyx/internal/utils"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Blocktunium/gonyx/internal/logger/helpers"
)

type OutputOption struct {
	LevelStr      string         `json:"level"`
	Level         types.LogLevel `json:"-"`
	Path          string         `json:"path,omitempty"`
	l             *log.Logger    `json:"-"`
	sqlDbInstance *gorm.DB       `json:"-"`
	dbType        string         `json:"-"`
}

// MARK: LogMeWrapper

// LogMeWrapper structure - implements Logger interface
type LogMeWrapper struct {
	name                  string
	serviceName           string
	ch                    chan types.LogObject
	initialized           bool
	wg                    sync.WaitGroup
	operationType         string
	supportedOutput       []string
	supportedOutputOption map[string]OutputOption
}

// Constructor - It initializes the logger configuration params
func (l *LogMeWrapper) Constructor(name string) error {
	l.wg.Add(1)
	defer l.wg.Done()

	l.name = name
	l.serviceName = config.GetManager().GetName()
	l.operationType = config.GetManager().GetOperationType()
	l.supportedOutput = []string{"console", "file", "db"}
	l.initialized = false

	channelSize, err := config.GetManager().Get(l.name, "channel_size")
	if err != nil {
		return err
	}

	options, err := config.GetManager().Get(l.name, "options")
	if err != nil {
		return err
	}

	optionArray := make([]string, len(options.([]interface{})))
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

	//if l.operationType == "prod" {
	//} else {
	//}
	l.supportedOutputOption = make(map[string]OutputOption)
	for _, item := range outputArray {
		if utils.ArrayContains(&l.supportedOutput, item) {
			jsonObj, configReadErr := config.GetManager().Get(l.name, item)
			if configReadErr == nil {
				r := OutputOption{}

				jsonMap := jsonObj.(map[string]interface{})
				jsonStr, jsonErr := json.Marshal(jsonMap)
				if jsonErr == nil {
					jsonErr = json.Unmarshal(jsonStr, &r)
					if jsonErr == nil {
						// add it to internal map
						r.Level = types.StringToLogLevel(r.LevelStr)
						if item == "console" {
							r.l = log.New(os.Stdout, "", 0)
						} else if item == "db" {
							useDbName := ""
							dbType := ""
							if v, ok := jsonMap["use"]; ok {
								useDbName = v.(string)
							}
							if v, ok := jsonMap["type"]; ok {
								dbType = v.(string)
							}

							if useDbName != "" && dbType != "" {
								if dbType == "sql" {
									r.dbType = dbType

									insDb, err := helpers.GetSqlDbInstance("server1")
									if err != nil {
										log.Printf("Cannot Db instance: %v for logger", item)
									} else {
										r.sqlDbInstance = insDb
										if r.sqlDbInstance != nil {
											err := r.sqlDbInstance.AutoMigrate(&types.ZhycanLog{})
											if err != nil {
												log.Printf("Cannot migrate the `ZhycanLog` table")
											}
										}
									}
								}
							}
						}
						l.supportedOutputOption[item] = r

						// run the instance ...
						go l.runner(item)
					} else {
						log.Printf("Cannot create log instance for: %v - %v", item, jsonErr)
					}
				} else {
					log.Printf("Cannot create log instance for: %v - %v", item, jsonErr)
				}
			} else {
				log.Printf("Cannot create log instance for: %v - %v", item, configReadErr)
			}
		} else {
			log.Printf("Log outout with name `%v` is not supported yet", item)
		}
	}

	l.initialized = true

	return nil
}

// IsInitialized - that returns boolean value whether it's initialized
func (l *LogMeWrapper) IsInitialized() bool {
	return l.initialized
}

// Instance - returns exact logger instance
func (l *LogMeWrapper) Instance() *log.Logger {
	l.wg.Wait()
	return l.supportedOutputOption["console"].l
}

// Log - write log object to the channel
func (l *LogMeWrapper) Log(obj *types.LogObject) {
	l.wg.Wait()

	go func(obj *types.LogObject) {
		l.ch <- *obj
	}(obj)
}

// Sync - sync all logs to medium
func (l *LogMeWrapper) Sync() {
	l.wg.Wait()
	ch := make(chan bool, 1)
	go func() {
		for {
			if len(l.ch) > 0 {
				time.Sleep(time.Millisecond * 200)
			} else {
				ch <- true
			}
		}
	}()

	<-ch
}

// Close - it closes logger channel
func (l *LogMeWrapper) Close() {
	l.wg.Wait()
	l.Sync()
	defer close(l.ch)
}

// runner - the goroutine that reads from channel and process it
func (l *LogMeWrapper) runner(output string) {
	l.wg.Wait()
	for c := range l.ch {
		if c.Level <= l.supportedOutputOption[output].Level {
			if output == "console" {
				switch c.Level {
				case types.DEBUG:
					l.debug(&c, output)
					break
				case types.INFO:
					l.info(&c, output)
					break
				case types.WARNING:
					l.warning(&c, output)
					break
				case types.ERROR:
					l.error(&c, output)
					break
				}
			} else if output == "db" {
				if l.supportedOutputOption[output].dbType == "sql" {
					if l.supportedOutputOption[output].sqlDbInstance != nil {
						item := types.ZhycanLog{
							Model:       gorm.Model{},
							ServiceName: l.serviceName,
							Level:       c.Level.String(),
							LogType:     c.LogType,
							Module:      c.Module,
							Message:     fmt.Sprintf("%v", c.Message),
							Additional:  fmt.Sprintf("%v", c.Additional),
							LogTime:     c.Time,
						}
						l.supportedOutputOption[output].sqlDbInstance.Create(&item)
					}
				}
			}
		}
	}
}

// debug - log with DEBUG level
func (l *LogMeWrapper) debug(object *types.LogObject, output string) {
	if output == "console" {
		l.supportedOutputOption[output].l.Printf(
			"\033[37m%v %v >>> %7v >>> (%v/%v)  - %v ... %v\033[0m\n",
			l.serviceName,
			object.Time,
			object.Level.String(),
			object.LogType,
			object.Module,
			object.Message,
			object.Additional,
		)
	}
}

// info - log with INFO level
func (l *LogMeWrapper) info(object *types.LogObject, output string) {
	if output == "console" {
		l.supportedOutputOption[output].l.Printf(
			"\033[32m%v %v >>> %7v >>> (%v/%v)  - %v ... %v\033[0m\n",
			l.serviceName,
			object.Time,
			object.Level.String(),
			object.LogType,
			object.Module,
			object.Message,
			object.Additional,
		)
	}
}

// warning - log with WARNING level
func (l *LogMeWrapper) warning(object *types.LogObject, output string) {
	if output == "console" {
		l.supportedOutputOption[output].l.Printf(
			"\033[33m%v %v >>> %7v >>> (%v/%v)  - %v ... %v\033[0m\n",
			l.serviceName,
			object.Time,
			object.Level.String(),
			object.LogType,
			object.Module,
			object.Message,
			object.Additional,
		)
	}
}

// error - log with ERROR level
func (l *LogMeWrapper) error(object *types.LogObject, output string) {
	if output == "console" {
		l.supportedOutputOption[output].l.Printf(
			"\033[31m%v %v >>> %7v >>> (%v/%v)  - %v ... %v\033[0m\n",
			l.serviceName,
			object.Time,
			object.Level.String(),
			object.LogType,
			object.Module,
			object.Message,
			object.Additional,
		)
	}
}
