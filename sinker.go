package main

import (
	"log"
	"os"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/gist"
)

func doIt(config conf.Conf, file conf.File) *gist.SyncResponse {
	fh, err := os.Open(file.Path)
	if err != nil {
		log.Fatalf("problem reading file: %s", err)
	}
	resp := gist.Sync(config.Gist.AccessToken, fh, file.Id)
	return &resp
}

func main() {

	config, err := conf.Get("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range config.Gist.Files {

		resp := doIt(*config, file)
		if resp.Error != nil {
			log.Fatal(resp.Error)
		}
		log.Println(resp)

	}
}
