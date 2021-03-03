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
	fmt.Println("get")
	return client(accessToken).Gists.Get(context.Background(), id)
}

func Sync(accessToken string, fh *os.File, gistId string) (string, error) {
	stat, err := fh.Stat()
	if err != nil {
		return "", fmt.Errorf("%w; could not get file stat", err)
	}
	fileUpdatedAt := stat.ModTime()
	gist, resp, err := Get(accessToken, gistId)
	if err != nil {
		return "", fmt.Errorf("couldn't get gist; %w", err)
	}
	if resp.Response.StatusCode != 200 {
		return "", fmt.Errorf("response from github was %d", resp.Response.StatusCode)
	}
	log.Printf("file %s last modified: %v\n", fh.Name(), fileUpdatedAt)
	log.Printf("gist last modified: %v\n", gist.UpdatedAt)
	log.Printf("file was modified after gist? %t\n", fileUpdatedAt.After(*gist.UpdatedAt))

	name := github.GistFilename(stat.Name())
	content := string(*gist.Files[name].Content)
	return content, nil
}
