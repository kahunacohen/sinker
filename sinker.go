package main

import (
	"conf/conf"
	"fmt"
	"log"
)

func main() {
	config, err := conf.Get()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(config.Gist.AccessToken)

}
