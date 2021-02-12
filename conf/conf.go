package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Gist struct {
	AccessToken string
	Files       []string
}

type Conf struct {
	Gist Gist
}

// Attempts to read  the .sinkerrc.json file in the user's
// home directory
func ReadSinkerRc() ([]byte, error) {
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

// Parses the json from the config.
func ParseJsonConfg(data []byte) (Conf, error) {
	var conf Conf
	err := json.Unmarshal(data, &conf)
	return conf, err
}

// Gets the configuration data as a Conf struct.
// The caller can directly reference fields on the struct
// because golang allows (*P).f to be accessed as P.f.
func Get() (*Conf, error) {
	data, err := ReadSinkerRc()
	if err != nil {
		fmt.Println("here")
		return nil, err
	}
	conf, err := ParseJsonConfg(data)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
