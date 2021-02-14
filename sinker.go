package main

import (
	"fmt"
	"log"

	"github.com/kahunacohen/sinker/conf"
	"github.com/kahunacohen/sinker/gist"
)

func main() {
	config, err := conf.Get("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(config.Gist.AccessToken)
	fmt.Println(gist.GetInfo("142a4dfb66f0e2eab38cb68e0b69d95c"))

}
