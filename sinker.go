package main

import (
	"fmt"
	"github.com/kahunacohen/sinker/conf"
	"log"
)

func main() {
	data, err := conf.ReadSinkerRc()
	if err != nil {
		log.Fatal("Problem reading your .sinkerrc.json file: " + err.Error())
	}
	config, err := conf.ParseJsonConfg(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config.Gist.Files[0])
	fmt.Println(config.Gist.AccessToken)

}
