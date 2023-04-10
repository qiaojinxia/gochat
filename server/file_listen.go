package server

import (
	"github.com/fsnotify/fsnotify"
	"goChat/conf"
	"log"
)

var fileListener *FileListener

type FileListener struct {
	done chan struct{}
}

func NewFileListener() *FileListener {
	return &FileListener{done: make(chan struct{})}
}

func (f *FileListener) Close() error {
	close(f.done)
	return nil
}

func (f *FileListener) Init() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer watcher.Close()
	err = watcher.Add(".")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				conf.InitGlobalConfig("config.ini")
			}
		case err := <-watcher.Errors:
			log.Println("Error: ", err)
		case <-f.done:
			log.Println("Exiting...")
			return nil
		}
	}
}
