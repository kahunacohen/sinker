package conf

import (
	"encoding/json"
	"io/ioutil"
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
func ReadSinkerRc(dir string) ([]byte, error) {
	data, err := ioutil.ReadFile(dir)
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
func Get(dir string) (*Conf, error) {
	data, err := ReadSinkerRc(dir)
	if err != nil {
		return nil, err
	}
	conf, err := ParseJsonConfg(data)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
