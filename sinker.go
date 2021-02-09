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
type Config struct {
	Gist
}

type Conf struct {
	Gist Gist
}

// Parses the json from the config.
func parseJsonConfg(data []byte) (Config, error) {
	var config Config
	err := json.Unmarshal([]byte(data), &config)
	return config, err
}
func main() {
	k := []byte(`{"gist": {"accessToken": "xxx", "files": ["~/.bashrc"]}}`)
	var c Conf
	err := json.Unmarshal(k, &c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.Gist.AccessToken)
	fmt.Println(c.Gist.Files)
	data, err := readSinkerRc()
	fmt.Println(string(data))
	if err != nil {
		log.Fatal("Problem reading your .sinkerrc.json file: " + err.Error())
	}
	config, err := parseJsonConfg(data)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(config)
	}

}
