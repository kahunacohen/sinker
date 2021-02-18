package gist

import (
	"log"

	"golang.org/x/oauth2"

	"context"

	"github.com/google/go-github/v33/github"
)

const gistApiUrl string = "https://api.github.com/gists/%s"

var c *github.Client = nil

func Client(accessToken string) *github.Client {
	if c == nil {
		log.Println("create auth client")
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		c = github.NewClient(tc)
	}
	return c
}

func Get(accessToken string, id string) (*github.Gist, *github.Response, error) {
	return Client(accessToken).Gists.Get(context.Background(), id)
}
