package config

// Imports needed list
import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Mark: manager

// Manager object
type manager struct {
	modules       map[string]*ViperWrapper
	modulesStatus map[string]bool

	configBasePath string
	configMode     string

	configRemoteAddress  string
	configRemoteInfra    string
	configRemoteDuration int64

	quitCh chan bool
}

// MARK: Module variables
var providerInstance *manager = nil
var once sync.Once

// MARK: Module Initializer
func init() {
	log.Println("Initializing Config Provider ...")
}

// MARK: Private Methods

// constructor - Constructor -> It initializes the config configuration params
func constructor(configBasePath string, configInitialMode string, configEnvPrefix string) error {
	log.Println("Config Manager Initializer ...")

	providerInstance.configMode = configInitialMode
	providerInstance.configBasePath = configBasePath

	viper.SetEnvPrefix(configEnvPrefix)
	err := viper.BindEnv("mode")
	if err != nil {
		return err
	}

	err = viper.BindEnv("name")
	if err != nil {
		return err
	}

	err = viper.BindEnv("config_remote_addr")
	if err != nil {
		return err
	}

	err = viper.BindEnv("config_remote_infra")
	if err != nil {
		return err
	}

	err = viper.BindEnv("config_remote_duration")
	if err != nil {
		return err
	}

	mode := viper.Get("mode")
	if mode != nil {
		providerInstance.configMode = mode.(string)
	}

	viper.AddConfigPath(fmt.Sprintf("%s/configs/%s/", providerInstance.configBasePath, providerInstance.configMode))
	viper.SetConfigName("base")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Load all modules
	configRemoteAddr := viper.GetString("config_remote_addr")
	configRemoteInfra := viper.GetString("config_remote_infra")
	configRemoteDuration := viper.GetInt64("config_remote_duration")

	providerInstance.configRemoteInfra = configRemoteInfra
	providerInstance.configRemoteAddress = configRemoteAddr
	providerInstance.configRemoteDuration = configRemoteDuration

	providerInstance.loadModules()

	log.Printf("Read Base `%s` Configs", viper.GetString("name"))
	mustWatched := viper.GetBool("config_must_watched")
	if mustWatched {
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Println("Configs Changed: ", in.Name)
		})
	}
	return nil
}

// loadModules - Loads All Modules That is configured in "init" config file
func (p *manager) loadModules() {
	log.Println("Load All Modules Config ...")
	modules := viper.Get("modules")

	for _, item2 := range modules.([]interface{}) {
		item := item2.(map[string]interface{})
		name := item["name"].(string)

		w := &ViperWrapper{
			ConfigPath:          []string{fmt.Sprintf("%s/configs/%s/", p.configBasePath, p.configMode)},
			ConfigName:          item["name"].(string),
			ConfigResourcePlace: item["type"].(string),
		}

		err := w.Load()
		if err == nil {
			p.modules[name] = w
			p.modulesStatus[name] = true
		} else {
			p.modulesStatus[name] = false
		}
	}
	//else if src == "remote" {
	//	// Start a goroutine
	//	for _, item := range modules {
	//		w := &ViperWrapper{
	//			ConfigName: item,
	//		}
	//
	//		p.modules[item] = w
	//		p.modulesStatus[item] = false
	//	}
	//
	//	// start remote loader as go routines
	//	go p.remoteConfigLoader()
	//}
}

// MARK: Public Methods

// CreateManager - Create a new manager instance
func CreateManager(configBasePath string, configInitialMode string, configEnvPrefix string) error {
	// once used for prevent race condition and manage critical section.
	if providerInstance == nil {
		var err error
		once.Do(func() {
			providerInstance = &manager{
				modules:       make(map[string]*ViperWrapper),
				modulesStatus: make(map[string]bool),
				quitCh:        make(chan bool),
			}

			for item := range providerInstance.modulesStatus {
				providerInstance.modulesStatus[item] = false
			}

			err = constructor(configBasePath, configInitialMode, configEnvPrefix)
		})
		return err
	}
	return nil
}

func GetManager() *manager {
	return providerInstance
}

// GetConfigWrapper - returns Config Wrapper based on name
func (p *manager) GetConfigWrapper(category string) (*ViperWrapper, error) {
	if val, ok := p.modules[category]; ok {
		return val, nil
	}

	return nil, NewCategoryNotExistErr(category, nil)
}

// GetName - returns service instance name based on config
func (p *manager) GetName() string {
	return viper.GetString("name")
}

// GetVersion - returns service instance name based on config
func (p *manager) GetVersion() string {
	return viper.GetString("version")
}

// GetOperationType - returns operation type which could be `dev`, `prod`
func (p *manager) GetOperationType() string {
	return p.configMode
}

// GetHostName - returns hostname based on config
func (p *manager) GetHostName() string {
	return os.Getenv(fmt.Sprintf("%s_HOSTNAME", p.GetName()))
}

// Get - get value of the key in specific category
func (p *manager) Get(category string, name string) (interface{}, error) {
	if val, ok := p.modules[category]; ok {
		result, exist := val.Get(name, false)
		if exist {
			return result, nil
		}

		return nil, NewKeyNotExistErr(name, category, nil)
	}

	return nil, NewCategoryNotExistErr(name, nil)
}

// Set - set value in category by specified key.
func (p *manager) Set(category string, name string, value interface{}) error {
	if val, ok := p.modules[category]; ok {
		return val.Set(name, value, false)
	}

	return NewCategoryNotExistErr(category, nil)
}

// StopLoader - stop remote loader
func (p *manager) StopLoader() {
	//if p.ConfigSrc == "remote" {
	//	p.quitCh <- true
	//}
}

// IsInitialized - iterate over all config wrappers and see all initialised correctly
func (p *manager) IsInitialized() bool {
	flag := true
	for _, value := range p.modulesStatus {
		if value == false {
			flag = false
			break
		}
	}
	return flag
}

// GetAllInitializedModuleList - get list of names that initialized truly
func (p *manager) GetAllInitializedModuleList() []string {
	var result []string
	for key, val := range p.modulesStatus {
		if val {
			result = append(result, key)
		}
	}

	return result
}

// remoteConfigLoader - get configs from remote
func (p *manager) remoteConfigLoader() {
	for {
		select {
		case <-p.quitCh:
			return
		default:
			for key := range p.modulesStatus {
				data, err := p.remoteConfigLoad(key)
				if err == nil {
					err = p.modules[key].LoadFromRemote(data)
					if err == nil {
						p.modulesStatus[key] = true
					} else {
						log.Println(err.Error())
					}
				} else {
					log.Println(err.Error())
				}
			}
		}

		time.Sleep(time.Duration(p.configRemoteDuration) * time.Second)
	}
}

// remoteConfigLoad
func (p *manager) remoteConfigLoad(key string) ([]byte, error) {
	//if p.ConfigRemoteAddress != "" {
	//	if p.ConfigRemoteInfra == "grpc" {
	//		conn, err := grpc.Dial(p.ConfigRemoteAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//		if err != nil {
	//			return nil, protoerror.GrpcDialError{Addr: p.ConfigRemoteAddress, Err: err}
	//		}
	//		defer conn.Close()
	//
	//		localContext, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	//		defer cancel()
	//
	//		c := api.NewHelpClient(conn)
	//		response, err := c.GetServiceConfig(localContext, &api.ServiceConfigRequest{
	//			Section:       key,
	//			Hostname:      p.GetHostName(),
	//			ConstructorId: api.ConstructorId_V1901,
	//			ServiceName:   p.GetName(),
	//			MsgId:         timestamppb.Now()})
	//		if err != nil {
	//			return nil, RemoteResponseError{Err: err}
	//		}
	//
	//		return []byte(response.Data), nil
	//	}
	//}
	return nil, NewRemoteLoadErr(key, nil)
}

// ManualLoadConfig - load manual config from the path and add to the current dict
func (p *manager) ManualLoadConfig(configBasePath string, configName string) error {
	w := &ViperWrapper{
		ConfigPath:          []string{configBasePath},
		ConfigName:          configName,
		ConfigResourcePlace: "",
	}

	err := w.Load()
	if err == nil {
		p.modules[configName] = w
		p.modulesStatus[configName] = true
	} else {
		p.modulesStatus[configName] = false
		return err
	}
	return nil
}
