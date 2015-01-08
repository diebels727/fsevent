package main

import (
	"flag"
	"github.com/diebels727/fsnotify"
	"log"
	"os/exec"
)

var watch string
var command string
var commandMessage string

func init() {
	flag.StringVar(&watch, "watch", "/tmp", "path or file to watch for changes")
	flag.StringVar(&command, "command", "", "OS command to execute; should be full path")
	flag.StringVar(&commandMessage, "message", "Executing.", "Message to display when this command executes.")
}

func main() {
	flag.Parse()
	log.Printf("Watching %s", watch)
	log.Printf("Command \"%s\" will execute on file change operations.", command)
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
					log.Println(commandMessage)
					cmd := exec.Command(command)
					err := cmd.Run()
					if err != nil {
						log.Println(err)
					}
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
