package gist

import (
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
)

const gistApiUrl string = "https://api.github.com/gists/%s"

// Gets a gist info by ID
func GetInfo(id string) (string, error) {
	res, err := http.Get(fmt.Sprintf(gistApiUrl, id))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code not 200 getting gist: %d", res.StatusCode)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
