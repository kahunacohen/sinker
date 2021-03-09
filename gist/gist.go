package gist

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"

	"context"

	"github.com/google/go-github/v33/github"
)

var c *github.Client = nil

// Wraps the github golang sdk authorized client.
func client(accessToken string) *github.Client {
	if c == nil {
		// log.Println("create auth client")
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

// A response from syncing to github.
// Content is the string representing the content
// of the file that should be synced, whether from
// the local file or from the remote gist.
// If the file is nil, that means the content represents
// the remote gist.
type SyncResponse struct {
	GistContent string
	File        *os.File
	FileNewer   bool
	Error       error
}

// Given a file handle and a gist ID returns a struct with the data needed
// to sync.
func Sync(accessToken string, fh *os.File, gistId string) SyncResponse {
	stat, err := fh.Stat()
	if err != nil {
		return SyncResponse{GistContent: "", File: nil, FileNewer: false, Error: err}
	}
	fileUpdatedAt := stat.ModTime()
	gist, resp, err := Get(accessToken, gistId)
	if err != nil {
		return SyncResponse{GistContent: "", File: nil, FileNewer: false, Error: err}
	}
	if resp.Response.StatusCode != 200 {
		return SyncResponse{GistContent: "", File: nil, Error: fmt.Errorf("response from github was %d", resp.Response.StatusCode)}
	}
	name := github.GistFilename(stat.Name())
	return SyncResponse{
		File:        fh,
		FileNewer:   fileUpdatedAt.After(*gist.UpdatedAt),
		GistContent: string(*gist.Files[name].Content),
		Error:       nil}
}
