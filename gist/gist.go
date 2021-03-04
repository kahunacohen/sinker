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

type SyncResponse struct {
	Local   bool
	Content string
	Error   error
}

func Sync(accessToken string, fh *os.File, gistId string) SyncResponse {
	// Given a file handle and a gist ID returns a struct indicating whether
	// The file or the gist's content is newer, the actual content and a possible
	// error.
	stat, err := fh.Stat()
	if err != nil {
		return SyncResponse{Local: false, Content: "", Error: err}
	}
	fileUpdatedAt := stat.ModTime()
	gist, resp, err := Get(accessToken, gistId)
	if err != nil {
		return SyncResponse{Local: false, Content: "", Error: err}
	}
	if resp.Response.StatusCode != 200 {
		return SyncResponse{Local: false, Content: "", Error: fmt.Errorf("response from github was %d", resp.Response.StatusCode)}
	}
	log.Printf("file %s last modified: %v\n", fh.Name(), fileUpdatedAt)
	log.Printf("gist last modified: %v\n", gist.UpdatedAt)
	localUpdatedLast := fileUpdatedAt.After(*gist.UpdatedAt)
	log.Printf("file was modified after gist? %t\n", fileUpdatedAt.After(*gist.UpdatedAt))
	name := github.GistFilename(stat.Name())
	return SyncResponse{Local: localUpdatedLast, Content: string(*gist.Files[name].Content), Error: nil}
}
