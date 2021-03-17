package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/gist"
)

func getOpts() conf.Opts {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "sinker syncs a set of local files to associated, remote gists, given a .sinker.json config file. Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	var verboseVar bool
	flag.BoolVar(&verboseVar, "verbose", false, "print log messages.")
	flag.Parse()
	return conf.Opts{Verbose: verboseVar}
}
func logResult(f conf.File, result gist.SyncResult) {
	log.Printf("syncing %s", f.Path)
	if result.Error != nil {
		log.Fatalf("sinker exited. Error: %v", result.Error)
	}
	if result.FileOverwritesGist {
		log.Println("file is newer, overwrote gist")
	} else if result.GistOverwritesFile {
		log.Println("gist newer, overwrote file")
	} else {
		log.Println("file and gist have the same content. noop")
	}
	fmt.Println()
}
func drainResultChannel(files []conf.File, syncResultChan chan gist.SyncResult, verbose bool) {
	for i, f := range files {
		result := <-syncResultChan
		if verbose {
			logResult(f, result)
		}
		if i == len(files)-1 {
			close(syncResultChan)
		}
	}
}
func main() {
	opts := getOpts()
	config, err := conf.Load("/Users/acohen/.sinkerrc.json", opts)
	if err != nil {
		log.Fatalf("unable to parse .sinkerrc.json: %v", err)
	}
	syncDataChan := make(chan gist.SyncData)
	syncResultChan := make(chan gist.SyncResult)
	for _, gistFile := range config.Gist.Files {
		go gist.GetSyncData(gistFile, syncDataChan, config)
		go gist.Sync(syncDataChan, syncResultChan, config)
	}
	drainResultChannel(config.Gist.Files, syncResultChan, opts.Verbose)
}
