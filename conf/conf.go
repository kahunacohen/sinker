package conf

import (
	"encoding/json"
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
