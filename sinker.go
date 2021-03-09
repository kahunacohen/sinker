package main

import (
	"fmt"
	"log"

	"path/filepath"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/compare"
	"github.com/kahunacohen/sinker/gist"
)

func main() {

	config, err := conf.Load("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err)
	}

	comparisonChan := make(chan *gist.SyncResponse, len(config.Gist.Files))
	for _, file := range config.Gist.Files {
		go compare.Compare(*config, file, comparisonChan)
	}
	for i := range config.Gist.Files {
		comparison := <-comparisonChan
		log.Printf("%s:", filepath.Base(comparison.File.Name()))
		if comparison.FileNewer {
			log.Println("The FILE is newer")
		} else {
			log.Println("The GIST is newer")
		}
		fmt.Println("")
		if i == len(config.Gist.Files)-1 {
			close(comparisonChan)
		}
	}
}
