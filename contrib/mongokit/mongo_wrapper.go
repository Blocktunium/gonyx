package mongokit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Blocktunium/gonyx/contrib/mongokit/extensions"
	"github.com/Blocktunium/gonyx/pkg/config"
	"github.com/Blocktunium/gonyx/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

// Mark: Definitions

// SqlWrapper struct
type MongoWrapper struct {
	name             string
	config           *Mongo
	databaseInstance *mongo.Client
}

func (m *MongoWrapper) init(name string) error {
	m.name = name

	// reading config
	nameParts := strings.Split(m.name, "/")

	var tempConfig *Mongo
	tempConfigObj, err := config.Get(nameParts[0], nameParts[1])
	if err == nil {
		// first marshal
		configData, err := json.Marshal(tempConfigObj)
		if err == nil {
			_ = json.Unmarshal(configData, &tempConfig)
		}
	}
	m.config = tempConfig

	//dbNameKey := fmt.Sprintf("%s.%s", nameParts[1], "db")
	//dbNameStr, err := config.Get(nameParts[0], dbNameKey)
	//if err != nil {
	//	return err
	//}
	//
	//hostKey := fmt.Sprintf("%s.%s", nameParts[1], "host")
	//hostStr, err := config.Get(nameParts[0], hostKey)
	//if err != nil {
	//	return err
	//}
	//
	//portKey := fmt.Sprintf("%s.%s", nameParts[1], "port")
	//portStr, err := config.Get(nameParts[0], portKey)
	//if err != nil {
	//	return err
	//}
	//
	//usernameKey := fmt.Sprintf("%s.%s", nameParts[1], "username")
	//usernameStr, err := config.Get(nameParts[0], usernameKey)
	//if err != nil {
	//	return err
	//}
	//
	//passwordKey := fmt.Sprintf("%s.%s", nameParts[1], "password")
	//passwordStr, err := config.Get(nameParts[0], passwordKey)
	//if err != nil {
	//	return err
	//}

	optionsKey := fmt.Sprintf("%s.%s", nameParts[1], "options")
	optionsObj, err := config.Get(nameParts[0], optionsKey)
	if err != nil {
		return err
	}

	optionsMap := make(map[string]string, len(optionsObj.(map[string]interface{})))
	for key, item := range optionsObj.(map[string]interface{}) {
		optionsMap[key] = item.(string)
	}
	m.config.Options = optionsMap

	var internalLogger *MongoLoggerConfig

	internalLoggerKey := fmt.Sprintf("%s.%s", nameParts[1], "logger")
	internalLoggerObj, err := config.Get(nameParts[0], internalLoggerKey)
	if err == nil {
		// first marshal
		configData, err := json.Marshal(internalLoggerObj)
		if err == nil {
			_ = json.Unmarshal(configData, &internalLogger)
		}
	}
	m.config.LoggerConfig = internalLogger

	//m.config = &Mongo{
	//	DatabaseName: dbNameStr.(string),
	//	Username:     usernameStr.(string),
	//	Password:     passwordStr.(string),
	//	Host:         hostStr.(string),
	//	Port:         portStr.(string),
	//	Options:      optionsMap,
	//}
	return nil
}

func (m *MongoWrapper) makeUri() string {
	optionsQSArr := make([]string, 0)
	for key, val := range m.config.Options {
		optionsQSArr = append(optionsQSArr, fmt.Sprintf("%s=%s", key, val))
	}
	optionsQS := strings.Join(optionsQSArr, "&")

	if m.config.Username == "" || m.config.Password == "" {
		return fmt.Sprintf("mongodb://%s:%s/%s?%s",
			m.config.Host,
			m.config.Port,
			m.config.DatabaseName,
			optionsQS)

	}

	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?%s",
		m.config.Username,
		m.config.Password,
		m.config.Host,
		m.config.Port,
		m.config.DatabaseName,
		optionsQS)
}

// MARK: Public functions

// GetDb - return associated internal Db
func (m *MongoWrapper) GetDb() (*mongo.Database, error) {
	if m.databaseInstance == nil {
		uri := m.makeUri()

		// Configure Logger Options
		loggerOptions := options.Logger().SetSink(&extensions.MongoLogger{}).SetMaxDocumentLength(uint(m.config.LoggerConfig.MaxDocumentLength))
		switch m.config.LoggerConfig.ComponentConnection {
		case "info":
			loggerOptions = loggerOptions.SetComponentLevel(options.LogComponentConnection, options.LogLevelInfo)
		default:
			loggerOptions = loggerOptions.SetComponentLevel(options.LogComponentConnection, options.LogLevelDebug)
		}

		switch m.config.LoggerConfig.ComponentCommand {
		case "info":
			loggerOptions = loggerOptions.SetComponentLevel(options.LogComponentCommand, options.LogLevelInfo)
		default:
			loggerOptions = loggerOptions.SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
		}

		db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetLoggerOptions(loggerOptions))
		if err != nil {
			return nil, err
		}

		ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
		errPing := db.Ping(ctx, nil)
		if errPing != nil {
			return nil, errPing
		}

		m.databaseInstance = db
	}

	actualDb := m.databaseInstance.Database(m.config.DatabaseName, nil)
	return actualDb, nil
}

// NewMongoWrapper - create a new instance of MongoWrapper and returns it
func NewMongoWrapper(name string) (*MongoWrapper, error) {
	wrapper := &MongoWrapper{}
	err := wrapper.init(name)
	if err != nil {
		return nil, NewCreateMongoWrapperErr(err)
	}

	return wrapper, nil
}

// RegisterLogger - register logger instance
func (m *MongoWrapper) RegisterLogger(l logger.Logger) {
	// MongoDB driver handles logging internally, so we store the logger for potential future use
	// The actual logging is handled by the mongo driver's options.LoggerOptions
}
