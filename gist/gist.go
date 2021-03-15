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

// Updates a gist
func UpdateGist(accessToken string, gist *github.Gist, gistFilename github.GistFilename, content []byte) (*github.Gist, *github.Response, error) {
	//GistContent: []byte(*gist.Files[fileNameFromGist].Content),
	// t := users[5]
	// t.name = "Mark"
	// users[5] = t

	// x := gist.Files[gistFilename]
	// ct := string(content)
	// x.Content = &ct
	// gist.Files[gistFilename] = x
	return client(accessToken).Gists.Edit(context.Background(), string(content), gist)
}

// A response from syncing to github.
// Content is the string representing the content
// of the file that should be synced, whether from
// the local file or from the remote gist.
// If the file is nil, that means the content represents
// the remote gist.
type SyncData struct {
	AccessToken  string
	Gist         *github.Gist
	GistFilename github.GistFilename
	FileContent  []byte
	FilePath     string
	FileLastMod  *time.Time

	Error error
}

// Given an access token, file handle, gist ID and a channel writes to channel a struct
// with the data needed to sync a local file and a gist.
func GetSyncData(accessToken string, gistFile conf.File, syncDataChan chan<- SyncData) {
	////fileNameFromGist := github.GistFilename(stat.Name())//fileNameFromGist := github.GistFilename(stat.Name())
	fh, err := os.Open(gistFile.Path)
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken:  accessToken,
			Gist:         nil,
			GistFilename: "",
			FileContent:  nil,
			FilePath:     gistFile.Path,
			FileLastMod:  nil,
			Error:        err}
		return
	}
	defer fh.Close()
	stat, err := fh.Stat()
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken:  accessToken,
			Gist:         nil,
			GistFilename: "",
			FileContent:  nil,
			FilePath:     gistFile.Path,
			FileLastMod:  nil,
			Error:        err}
		return
	}
	fileLastMod := stat.ModTime()
	gist, resp, err := GetGistData(accessToken, gistFile.Id)
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken:  accessToken,
			Gist:         nil,
			GistFilename: "",
			FileContent:  nil,
			FilePath:     gistFile.Path,
			FileLastMod:  &fileLastMod,
			Error:        err}
		return
	}
	if resp.Response.StatusCode != 200 {
		syncDataChan <- SyncData{
			AccessToken:  accessToken,
			Gist:         nil,
			GistFilename: "",
			FileContent:  nil,
			FilePath:     gistFile.Path,
			FileLastMod:  &fileLastMod,
			Error:        fmt.Errorf("response from github was %d")}
	}

	// Get the filename from gist so we can index into the files map.
	//fileNameFromGist := github.GistFilename(stat.Name())
	fileContent, err := ioutil.ReadFile(gistFile.Path)
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken:  accessToken,
			Gist:         nil,
			GistFilename: "",
			FileContent:  nil,
			FilePath:     gistFile.Path,
			FileLastMod:  &fileLastMod,
			Error:        err}
		return
	}
	syncDataChan <- SyncData{
		AccessToken:  accessToken,
		Gist:         gist,
		GistFilename: github.GistFilename(stat.Name()),
		FileContent:  fileContent,
		FilePath:     gistFile.Path,
		FileLastMod:  &fileLastMod,
		Error:        nil}
}

func Sync(syncDataChan <-chan SyncData, syncChan chan<- bool) {
	data := <-syncDataChan
	gist := *data.Gist
	log.Printf("syncing %s", data.FilePath)
	gistContent := gist.Files[data.GistFilename].Content
	if *gistContent == string(data.FileContent) {
		log.Printf("content is equal for file and gist--Do nothing.")
		syncChan <- true
		return
	}
	if data.FileLastMod.After(gist.UpdatedAt) {
		log.Println("the file is newer, push file contents to gist")
		_, resp, err := UpdateGist(data.AccessToken, gist, data.GistFilename, data.FileContent)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp.StatusCode)

	} else {
		log.Printf("the gist is newer, overwrite file %s", data.FilePath)
		gistContent := []byte(*gist.Files[data.GistFilename].Content)
		err := ioutil.WriteFile(data.FilePath, gistContent, 0644)
		if err != nil {
			fmt.Println(err)
			syncChan <- true
			return
		}
	}
	syncChan <- true
}
