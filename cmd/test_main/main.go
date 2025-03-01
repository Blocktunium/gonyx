package main

import (
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Gonyx Project")

	path := "/Users/abolfazl.beh/Projects/Gonyx/"
	initialMode := "dev"
	prefix := "Gonyx"

	err := config.CreateManager(path, initialMode, prefix)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		flag := config.GetManager().IsInitialized()
		if flag {
			break
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("All modules is initialized")

	//l, _ := logger.GetManager().GetLogger()
	//if l != nil {
	//	l.Log(logger.NewLogObject(
	//		logger.INFO, "test", logger.FuncMaintenanceType, time.Now().UTC(), "this is a test", nil))
	//}

	AddWatcher()

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

}

func AddWatcher() {
	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Move, watcher.Remove, watcher.Rename, watcher.Write)
	w.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile(".*.go$"), false))
	if err := w.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}
	w.SetMaxEvents(1)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Print a list of all the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}

}
