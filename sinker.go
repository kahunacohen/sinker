package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/gist"
)

func doIt(config conf.Conf, file conf.File) {
	fh, err := os.Open(file.Path)
	if err != nil {
		log.Fatalf("problem reading file: %s", err)
	}
	resp := gist.Sync(config.Gist.AccessToken, fh, file.Id)
	if resp.Error != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(resp.LocalModLast)
}

func main() {

	config, err := conf.Get("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range config.Gist.Files {
		doIt(*config, file)

	}
}
