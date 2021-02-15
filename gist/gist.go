package gist

import (
	"encoding/json"
	"log"

	"golang.org/x/oauth2"

	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/v33/github"
)

const gistApiUrl string = "https://api.github.com/gists/%s"

var _c *github.Client = nil

func Client(accessToken string) *github.Client {
	if _c == nil {
		log.Println("create auth client")
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		_c = github.NewClient(tc)
	}
	return _c
}

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
