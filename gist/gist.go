package gist

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"

	"context"

	"github.com/google/go-github/v33/github"
	"github.com/kahunacohen/sinker/conf"
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
func UpdateGist(accessToken string, id string, content []byte) (*github.Gist, *github.Response, error) {
	return client(accessToken).Gists.Edit(context.Background(), id, content)
}

// A response from syncing to github.
// Content is the string representing the content
// of the file that should be synced, whether from
// the local file or from the remote gist.
// If the file is nil, that means the content represents
// the remote gist.
type SyncData struct {
	AccessToken string
	GistContent []byte
	FileContent []byte
	FilePath    string
	GistId      string
	FileLastMod *time.Time
	GistLastMod *time.Time
	Error       error
}

// Given an access token, file handle, gist ID and a channel writes to channel a struct
// with the data needed to sync a local file and a gist.
func GetSyncData(accessToken string, gistFile conf.File, syncDataChan chan<- SyncData) {
	fh, err := os.Open(gistFile.Path)
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken: accessToken,
			GistContent: nil,
			FileContent: nil,
			GistId:      gistFile.Id,
			FilePath:    gistFile.Path,
			FileLastMod: nil,
			GistLastMod: nil,
			Error:       err}
		return
	}
	defer fh.Close()
	stat, err := fh.Stat()
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken: accessToken,
			GistContent: nil,
			FileContent: nil,
			GistId:      gistFile.Id,
			FilePath:    gistFile.Path,
			FileLastMod: nil,
			GistLastMod: nil,
			Error:       err}
		return
	}
	fileLastMod := stat.ModTime()
	gistData, resp, err := GetGistData(accessToken, gistFile.Id)
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken: accessToken,
			GistContent: nil,
			FileContent: nil,
			GistId:      gistFile.Id,
			FilePath:    gistFile.Path,
			FileLastMod: &fileLastMod,
			GistLastMod: nil,
			Error:       err}
		return
	}
	if resp.Response.StatusCode != 200 {
		syncDataChan <- SyncData{
			AccessToken: accessToken,
			GistContent: nil,
			FileContent: nil,
			GistId:      gistFile.Id,
			FilePath:    gistFile.Path,
			FileLastMod: &fileLastMod,
			GistLastMod: nil,
			Error:       fmt.Errorf("response from github was %d")}
	}

	// Get the filename from gist so we can index into the files map.
	fileNameFromGist := github.GistFilename(stat.Name())
	bytes, err := ioutil.ReadFile(gistFile.Path)
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken: accessToken,
			GistContent: nil,
			FileContent: nil,
			GistId:      gistFile.Id,
			FilePath:    gistFile.Path,
			FileLastMod: &fileLastMod,
			GistLastMod: nil,
			Error:       err}
		return
	}
	syncDataChan <- SyncData{
		AccessToken: accessToken,
		GistContent: []byte(*gistData.Files[fileNameFromGist].Content),
		FileContent: bytes,
		GistId:      gistFile.Id,
		FilePath:    gistFile.Path,
		FileLastMod: &fileLastMod,
		GistLastMod: gistData.UpdatedAt,
		Error:       nil}
}

func Sync(syncDataChan <-chan SyncData, syncChan chan<- bool) {
	data := <-syncDataChan
	log.Printf("syncing %s", data.FilePath)

	if string(data.GistContent) == string(data.FileContent) {
		log.Printf("content is equal for file and gist--Do nothing.")
		syncChan <- true
		return
	}
	if data.FileLastMod.After(*data.GistLastMod) {
		log.Println("the file is newer, push file contents to gist")
		_, _, err := UpdateGist(data.AccessToken, data.GistId, data.FileContent)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		log.Printf("the gist is newer, overwrite file %s", data.FilePath)
		err := ioutil.WriteFile(data.FilePath, data.GistContent, 0644)
		if err != nil {
			fmt.Println(err)
			syncChan <- true
			return
		}
	}
	syncChan <- true
}
