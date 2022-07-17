package framework

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

var done = make(chan bool)
var Reload = make(chan bool, 1)

func ReloadBrowser() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("watcher", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					Reload <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	gotmpls := GetFilesAndDirRecursively("static/pages")

	for _, tmpl := range gotmpls {
		err = watcher.Add(tmpl)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done
}

func CloseReloadBrowser() {
	done <- true
}
