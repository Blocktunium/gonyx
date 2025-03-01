package watcher

import (
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/radovskyb/watcher"
	"log"
	"regexp"
	"sync"
	"time"
)

// Mark: manager

// Manager object
type manager struct {
	name              string
	watcher           *watcher.Watcher
	lock              sync.Mutex
	isInitialized     bool
	printWatchedFiles bool
	watchInterval     int
}

type watchDirStruct struct {
	Path      string
	Recursive bool
}

// MARK: Module variables
var managerInstance *manager = nil
var once sync.Once

// MARK: Module Initializer
func init() {
	log.Println("Initializing Watcher Manager ...")
}

// MARK: Private Methods

// init - Constructor -> It initializes the config configuration params
func (m *manager) init() {
	log.Println("Constructing Watcher Manager ...")

	m.name = "watcher"
	m.isInitialized = false
	m.printWatchedFiles = false
	m.watchInterval = 100

	m.lock.Lock()
	defer m.lock.Unlock()

	// read filter options
	filterOptions, err := config.GetManager().Get(m.name, "filter_operations")
	if err != nil {
		return
	}

	var filterOptionsArray []string
	for _, v := range filterOptions.([]interface{}) {
		filterOptionsArray = append(filterOptionsArray, v.(string))
	}

	filterHooks, err := config.GetManager().Get(m.name, "filter_hooks")
	if err != nil {
		return
	}

	var filterHooksArray []string
	for _, v := range filterHooks.([]interface{}) {
		filterHooksArray = append(filterHooksArray, v.(string))
	}

	maxEvent, err := config.GetManager().Get(m.name, "max_event")
	if err != nil {
		return
	}

	printWatchedFiles, err := config.GetManager().Get(m.name, "print_watched_files")
	if err != nil {
		return
	}
	m.printWatchedFiles = printWatchedFiles.(bool)

	watchInterval, err := config.GetManager().Get(m.name, "watch_interval")
	if err != nil {
		return
	}
	m.watchInterval = watchInterval.(int)

	var watchDirs []watchDirStruct
	watchDirsObj, err := config.GetManager().Get(m.name, "watch_dirs")
	if err != nil {
		return
	}

	for _, v := range watchDirsObj.([]interface{}) {
		item := v.(map[string]interface{})
		watchDirs = append(watchDirs, watchDirStruct{
			Path:      item["path"].(string),
			Recursive: item["recursive"].(bool),
		})
	}

	// Let's config a watcher
	m.watcher = watcher.New()

	var filterOps []watcher.Op
	for _, item := range filterOptionsArray {
		switch item {
		case "create":
			filterOps = append(filterOps, watcher.Create)
			break
		case "move":
			filterOps = append(filterOps, watcher.Move)
			break
		case "rename":
			filterOps = append(filterOps, watcher.Rename)
			break
		case "remove":
			filterOps = append(filterOps, watcher.Remove)
			break
		case "write":
			filterOps = append(filterOps, watcher.Write)
			break
		}
	}

	m.watcher.FilterOps(filterOps...)

	for _, item := range filterHooksArray {
		m.watcher.AddFilterHook(watcher.RegexFilterHook(
			regexp.MustCompile(item), false))
	}

	m.watcher.SetMaxEvents(maxEvent.(int))

	for _, item := range watchDirs {
		if item.Recursive {
			if err := m.watcher.AddRecursive(item.Path); err != nil {
				continue
			}
		} else {
			if err := m.watcher.Add(item.Path); err != nil {
				continue
			}
		}
	}

	m.isInitialized = true
}

// restartOnChangeConfig - subscribe a function for when the config is changed
func (m *manager) restartOnChangeConfig() {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Config config server to reload
	wrapper, err := config.GetManager().GetConfigWrapper(m.name)
	if err == nil {
		wrapper.RegisterChangeCallback(func() interface{} {
			m.Stop()
			m.init()
			m.Start()
			return nil
		})
	} else {
		// TODO: make some logs
	}
}

// MARK: Public Functions

// GetManager - This function returns singleton instance of Logger Manager
func GetManager() *manager {
	// once used for prevent race condition and manage critical section.
	once.Do(func() {
		managerInstance = &manager{}
		managerInstance.init()
		managerInstance.restartOnChangeConfig()
	})
	return managerInstance
}

// Start - start the watcher and listen for changes
func (m *manager) Start() (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.printWatchedFiles {
		for path, f := range m.watcher.WatchedFiles() {
			fmt.Printf("%s: %s\n", path, f.Name())
		}
	}

	if err := m.watcher.Start(time.Millisecond * time.Duration(m.watchInterval)); err != nil {
		return false, NewStartWatcherErr(err)
	}

	return true, nil
}

// Stop - stop the watcher and listen for changes
func (m *manager) Stop() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.watcher.Close()
	<-m.watcher.Closed
}
