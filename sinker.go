package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kahunacohen/sinker/conf"
	"github.com/kahunacohen/sinker/gist"
)

func f1(ch chan int) {
	fmt.Println("f1 starting")
	time.Sleep(time.Second * 5)
	fmt.Println("f1 done")
	ch <- 1
}
func f2(ch chan int) {
	fmt.Println("f2 starting")
	time.Sleep(time.Second * 5)
	fmt.Println("f2 done")
	ch <- 2
}
func f3(ch chan int) {
	fmt.Println("f3 starting")
	time.Sleep(time.Second * 5)
	fmt.Println("f3 done")
	ch <- 3
}
func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)
	go f1(c1)
	go f2(c2)
	go f3(c3)
	x := <-c1
	y := <-c2
	z := <-c3
	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(z)

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
