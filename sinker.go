package main

import (
	"context"
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
	c := gist.Client(config.Gist.AccessToken)
	gist, resp, err := c.Gists.Get(context.Background(), "142a4dfb66f0e2eab38cb68e0b69d95c")
	if err != nil {
		log.Fatalf("couldn't get gist: %s", err)
	}

	if resp.Response.StatusCode != 200 {
		log.Fatalf("response from github was %d", resp.Response.StatusCode)
	}
	fmt.Println(gist)
}
