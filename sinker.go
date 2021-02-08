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

type AccessToken struct {
	AccessToken string `json:"access_token"`
}
type Entities struct {
	Entities []string `json:"entities"`
}
type Gist struct {
	AccessToken
	Entities
}

// Parses the json from the config.
func parseJsonConfg(data []byte) (Gist, error) {
	var gist Gist
	err := json.Unmarshal([]byte(data), &gist)
	return gist, err
}
func main() {
	data, err := readSinkerRc()
	if err != nil {
		log.Fatal("Problem reading your .sinkerrc.json file: " + err.Error())
	}
	fmt.Println(string(data))
	gist, err := parseJsonConfg(data)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(gist)
	}

}
