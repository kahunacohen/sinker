package gist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
)

const gistApiUrl string = "https://api.github.com/gists/%s"

// Gets a gist info by ID
func GetInfo(id string) (map[string]interface{}, error) {
	res, err := http.Get(fmt.Sprintf(gistApiUrl, id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code not 200 getting gist: %d", res.StatusCode)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var objmap map[string]interface{}
	err = json.Unmarshal(bytes, &objmap)
	if err != nil {
		return nil, err
	}
	return objmap, nil
}
