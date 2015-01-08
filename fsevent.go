package main

import (
	"flag"
	"github.com/diebels727/fsnotify"
	"log"
)

var watch string
var cmd string

func init() {
	flag.StringVar(&watch, "watch", "/tmp", "path or file to watch for changes")
	flag.StringVar(&cmd, "cmd", "", "OS command to execute")
}

func main() {
	flag.Parse()
	log.Printf("Watching %s", watch)
	log.Printf("  Command \"%s\" will execute on file change operations.", cmd)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				if event.Op&fsnotify.Close == fsnotify.Close {
					log.Println("Will fire!")
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(watch)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
