package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Attempts to read  the .sinkerrc.json file in the user's
// home directory
func readSinkerRc() ([]byte, error) {
	homdir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(path.Join(homdir, ".sinkerrc.json"))
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Gist struct {
	AccessToken string
	Files       []string
}

type Conf struct {
	Gist Gist
}

// Parses the json from the config.
func parseJsonConfg(data []byte) (Conf, error) {
	var conf Conf
	err := json.Unmarshal(data, &conf)
	return conf, err
}
func main() {
	data, err := readSinkerRc()
	if err != nil {
		log.Fatal("Problem reading your .sinkerrc.json file: " + err.Error())
	}
	conf, err := parseJsonConfg(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf.Gist.Files[0])
	fmt.Println(conf.Gist.AccessToken)

}
