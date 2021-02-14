package main

import (
	"fmt"
	"log"

	"github.com/kahunacohen/sinker/conf"
)

func main() {
	config, err := conf.Get("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(config.Gist.AccessToken)

}
