package main

import (
	"log"

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
		if comparison.File != nil {
			log.Println("The file is newer")
		} else {
			log.Println("The gist is newer")
		}
		if i == len(config.Gist.Files)-1 {
			close(comparisonChan)
		}
	}
}
