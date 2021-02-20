package main

import (
	"fmt"
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
		stat, err := fh.Stat()
		if err != nil {
			log.Fatal("could not get file stat: %s", err)
		}
		fileUpdatedAt := stat.ModTime()
		fmt.Printf("file last modified: %v\n", fileUpdatedAt)
		gist, resp, err := gist.Get(config.Gist.AccessToken, file.Id)
		if err != nil {
			log.Fatalf("couldn't get gist: %s", err)
		}
		if resp.Response.StatusCode != 200 {
			log.Fatalf("response from github was %d", resp.Response.StatusCode)
		}
		fmt.Printf("gist last modified: %v\n", gist.UpdatedAt)
		fmt.Printf("File was modified after gist? %t\n", fileUpdatedAt.After(*gist.UpdatedAt))
	}
}
