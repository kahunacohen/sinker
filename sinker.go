package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Attempts to read  the .sinkerrc.json file in the user's
// home directory
func readSinkerRc() (string, error) {
	homdir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadFile(path.Join(homdir, ".sinkerrc.json"))
	if err != nil {
		return "", err
	}
	return string(dat), nil
}
func main() {
	dat, err := readSinkerRc()
	if err != nil {
		log.Fatal("Problem reading your .sinkerrc.json file: " + err.Error())
	}
	fmt.Println(dat)
}
