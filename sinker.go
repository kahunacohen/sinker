package main

import (
	"fmt"
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
		fmt.Println(comparison.LocalModLast)
		if i == len(config.Gist.Files)-1 {
			close(comparisonChan)
		}
	}
}
