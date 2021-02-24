package gist

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"

	"context"

	"github.com/google/go-github/v33/github"
)

var c *github.Client = nil

// Wraps the github golang sdk authorized client.
func client(accessToken string) *github.Client {
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

// Get a gist named given a personal access token and a gist ID.
func Get(accessToken string, id string) (*github.Gist, *github.Response, error) {
	return client(accessToken).Gists.Get(context.Background(), id)
}

func Sync(accessToken string, fh *os.File, gistId string) ([]byte, error) {
	stat, err := fh.Stat()
	if err != nil {
		return nil, fmt.Errorf("%w; could not get file stat", err)
	}
	fileUpdatedAt := stat.ModTime()
	gist, resp, err := Get(accessToken, gistId)
	if err != nil {
		return nil, fmt.Errorf("couldn't get gist; %w", err)
	}
	if resp.Response.StatusCode != 200 {
		return nil, fmt.Errorf("response from github was %d", resp.Response.StatusCode)
	}
	log.Printf("file %s last modified: %v\n", fh.Name(), fileUpdatedAt)
	log.Printf("gist %s last modified: %v\n", "foo", gist.UpdatedAt)
	log.Printf("file was modified after gist? %t\n", fileUpdatedAt.After(*gist.UpdatedAt))

	return []byte("ABC"), nil
}
