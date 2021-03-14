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

// GetGistData a gist named given a personal access token and a gist ID.
func GetGistData(accessToken string, id string) (*github.Gist, *github.Response, error) {
	return client(accessToken).Gists.Get(context.Background(), id)
}

// A response from syncing to github.
// Content is the string representing the content
// of the file that should be synced, whether from
// the local file or from the remote gist.
// If the file is nil, that means the content represents
// the remote gist.
type SyncData struct {
	GistContent string
	File        *os.File
	FileNewer   bool
	Error       error
}

// Given an access token, file handle, gist ID and a channel writes to channel a struct
// with the data needed to sync a local file and a gist.
func GetSyncData(accessToken string, localFh *os.File, gistId string, syncDataChan chan<- SyncData) {
	stat, err := localFh.Stat()
	if err != nil {
		syncDataChan <- SyncData{GistContent: "", File: nil, FileNewer: false, Error: err}
	}
	fileUpdatedAt := stat.ModTime()
	gistData, resp, err := GetGistData(accessToken, gistId)
	if err != nil {
		syncDataChan <- SyncData{GistContent: "", File: nil, FileNewer: false, Error: err}
	}
	if resp.Response.StatusCode != 200 {
		syncDataChan <- SyncData{GistContent: "", File: nil, Error: fmt.Errorf("response from github was %d", resp.Response.StatusCode)}
	}

	// Get the filename from gist so we can index into the files map.
	fileNameFromGist := github.GistFilename(stat.Name())
	syncDataChan <- SyncData{
		File:        localFh,
		FileNewer:   fileUpdatedAt.After(*gistData.UpdatedAt),
		GistContent: string(*gistData.Files[fileNameFromGist].Content),
		Error:       nil}
}

func Sync(syncDataChan <-chan SyncData, syncChan chan<- bool) {
	data := <-syncDataChan
	syncChan <- data.FileNewer
}
