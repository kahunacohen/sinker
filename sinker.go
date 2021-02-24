package main

import (
	"log"
	"os"

	"github.com/kahunacohen/sinker/conf"
	"github.com/kahunacohen/sinker/gist"
)

func main() {

	config, err := conf.Get("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range config.Gist.Files {
		fh, err := os.Open(file.Path)
		if err != nil {
			log.Fatalf("problem reading file: %s", err)
		}
		_, err = gist.Sync(config.Gist.AccessToken, fh, file.Id)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
