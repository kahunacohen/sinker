package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/gist"
)

func getOpts() map[string]interface{} {
	ret := make(map[string]interface{})
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "sinker syncs a set of local files to associated, remote gists, given a .sinker.json config file. Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	var verboseVar bool
	flag.BoolVar(&verboseVar, "verbose", false, "print log messages.")
	flag.Parse()
	ret["verbose"] = verboseVar
	return ret
}
func main() {
	opts := getOpts()
	fmt.Println(opts)
	config, err := conf.Load("/Users/acohen/.sinkerrc.json")
	if err != nil {
		log.Fatal(err)
	}

	syncDataChan := make(chan gist.SyncData)
	syncChan := make(chan bool)
	for _, gistFile := range config.Gist.Files {
		go gist.GetSyncData(config.Gist.AccessToken, gistFile, syncDataChan)
		go gist.Sync(syncDataChan, syncChan)

	}

	for i := range config.Gist.Files {
		<-syncChan
		if i == len(config.Gist.Files)-1 {
			close(syncChan)
		}
	}

}
