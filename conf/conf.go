package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Opts struct {
	Verbose bool
}

type File struct {
	Path string
	Id   string
}
type Gist struct {
	AccessToken string
	Files       []File
}

type Conf struct {
	Gist Gist
	Opts Opts
}

// Attempts to read  the .sinkerrc.json file in the user's
// home directory
func readSinkerRc(dir string) ([]byte, error) {
	data, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Parses the json from the config.
func parseJSONConfig(data []byte) (Conf, error) {
	var conf Conf
	err := json.Unmarshal(data, &conf)
	return conf, err
}

// Load gets the configuration data as a Conf struct.
// The caller can directly reference fields on the struct
// because golang allows (*P).f to be accessed as P.f.
// The opts map is from the command line flags to be merged into the
// config.
func Load(dir string, opts Opts) (*Conf, error) {
	data, err := readSinkerRc(dir)
	if err != nil {
		return nil, err
	}
	conf, err := parseJSONConfig(data)
	conf.Opts = opts
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
