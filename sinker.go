package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/gist"
)

func main() {

	config, err := conf.Load("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err)
	}

	syncDataChan := make(chan *gist.SyncData, len(config.Gist.Files))
	for _, file := range config.Gist.Files {
		fh, err := os.Open(file.Path)
		if err != nil {
			log.Fatalf("problem reading file: %s", err)
		}
		go gist.GetSyncData(config.Gist.AccessToken, fh, file.Id, syncDataChan)

	}
	for i := range config.Gist.Files {
		syncData := <-syncDataChan
		log.Printf("%s:", filepath.Base(syncData.File.Name()))
		if syncData.FileNewer {
			log.Println("The FILE is newer")
		} else {
			log.Println("The GIST is newer")
		}
		fmt.Println("")
		if i == len(config.Gist.Files)-1 {
			close(syncDataChan)
		}
	}

}
