// A script that syncs gists. You must have a `.sinkerrc.json` file in your home directory in this form:
// {
// 	   "gist": {
// 				"accessToken": "xxx",
// 				"files": [
// 					{"path": "path/to/file/on/your/filesystem", "id": "gistid"}
// 				]
// 		 }
// }

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
	// Parse options from the command-line into a struct.
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
	// Since we know how many writes to the channel occured...it's one for each file,
	// we loop through the files, log the results, and close the channel once we're done.
	for i, f := range files {
		if verbose {
			result := <-syncResultChan
			logResult(f, result)
		}
		if i == len(files)-1 {
			close(syncResultChan)
		}
	}
}

// Entry point.
func main() {
	opts := getOpts()

	// Go handles errors by returning errors as values along with result.
	// This is not enforced by the type system, rather it's idiomatic.
	config, err := conf.Load("/Users/acohen/.sinkerrc.json", opts)
	// By convention, we check error conditions first before doing our main
	// operations. We prefix the existing error with more context usually as a string.
	if err != nil {
		log.Fatalf("unable to parse .sinkerrc.json: %v", err)
	}

	// Create 2 channels:
	// 1. Holds data for each file needed to sync.
	// 2. Holds results from actual syncing.
	syncDataChan := make(chan gist.SyncData)
	syncResultChan := make(chan gist.SyncResult)

	// Loop over the files defined in our config. For each one,
	// asyncronously get the data we need to sync, writing to the data channel and
	// pipe that channel to the function for syncing.
	for _, gistFile := range config.Gist.Files {
		go gist.GetSyncData(gistFile, syncDataChan, config)
		go gist.Sync(syncDataChan, syncResultChan, config)
	}

	// If we don't drain the result channel, the main function will exit before
	// our async functions complete. Draining, not only assures that our process
	// stays alive while workers do their work, here we also use it to report debug
	// data when we pass the --verbose flag.
	drainResultChannel(config.Gist.Files, syncResultChan, opts.Verbose)
}
