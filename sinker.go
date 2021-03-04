package main

import (
	"log"
	"os"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/gist"
)

func doIt(config conf.Conf, file conf.File, which chan *gist.SyncResponse) {
	fh, err := os.Open(file.Path)
	if err != nil {
		log.Fatalf("problem reading file: %s", err)
	}
	resp := gist.Sync(config.Gist.AccessToken, fh, file.Id)
	which <- &resp
}

func main() {

	config, err := conf.Get("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	which := make(chan *gist.SyncResponse, len(config.Gist.Files))
	for _, file := range config.Gist.Files {
		doIt(*config, file, which)
	}
	for i := range config.Gist.Files {
		<-which
		if i == len(config.Gist.Files)-1 {
			close(which)
		}
	}
}
