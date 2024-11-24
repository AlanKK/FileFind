package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsevents"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <root_folder>")
		os.Exit(1)
	}

	rootFolder := os.Args[1]
	absPath, err := filepath.Abs(rootFolder)
	if err != nil {
		log.Fatal(err)
	}

	dev, err := fsevents.DeviceForPath(absPath)
	if err != nil {
		log.Fatal(err)
	}

	es := &fsevents.EventStream{
		Paths:   []string{absPath},
		Latency: 500 * time.Millisecond,
		Device:  dev,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot,
	}

	es.Start()
	defer es.Stop()

	for {
		select {
		case events := <-es.Events:
			for _, event := range events {
				if event.Flags&fsevents.ItemCreated != 0 {
					printEvent("Created", event.Path)
				}
				if event.Flags&fsevents.ItemRemoved != 0 {
					printEvent("Deleted", event.Path)
				}
				if event.Flags&fsevents.ItemRenamed != 0 {
					printEvent("Renamed", event.Path)
				}
			}
		}
	}
}

func printEvent(eventType, path string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s: %s\n", timestamp, eventType, path)
}
