package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/gist"
)

func main() {
	// Loop through each file in config.
	// Get whether file or gist should be updated.

	config, err := conf.Load("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err)
	}

	syncDataChan := make(chan gist.SyncData)
	syncChan := make(chan bool)
	for _, file := range config.Gist.Files {
		fh, err := os.Open(file.Path)
		if err != nil {
			log.Fatalf("problem reading file: %s", err)
		}
		go gist.GetSyncData(config.Gist.AccessToken, fh, file.Id, syncDataChan)
		go gist.Sync(syncDataChan, syncChan)

	}

	for i := range config.Gist.Files {
		fmt.Println(<-syncChan)
		// log.Printf("%s:", filepath.Base(syncData.File.Name()))

		if i == len(config.Gist.Files)-1 {
			close(syncChan)
		}
	}

}
