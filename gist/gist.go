package gist

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/oauth2"

	"context"

	"github.com/google/go-github/v33/github"
	"github.com/kahunacohen/sinker/conf"
)

// Wraps the github golang sdk authorized client.
func client(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)

}

func getGistData(accessToken string, id string) (*github.Gist, *github.Response, error) {
	return client(accessToken).Gists.Get(context.Background(), id)
}

func updateGist(accessToken string, gist *github.Gist, gistFilename github.GistFilename, content []byte) (*github.Gist, *github.Response, error) {
	// We get a panic trying to index a map and update a struct field in place, so copy the
	// structure.
	x := gist.Files[gistFilename]
	ct := string(content)

	// Reassign the content to the gist in-place, then pass to Gists.Edit.
	x.Content = &ct

	gist.Files[gistFilename] = x
	return client(accessToken).Gists.Edit(context.Background(), *gist.ID, gist)
}

// SyncData is the data needed to sync a remote gist
// with a local file.
type SyncData struct {
	AccessToken  string
	Gist         *github.Gist
	GistFilename github.GistFilename
	FileContent  []byte
	FilePath     string
	FileLastMod  *time.Time
	Error        error
}

// SyncResult represents the result of syncing a local file and a remote Gist.
type SyncResult struct {
	FileOverwritesGist bool
	GistOverwritesFile bool
	Error              error
}

// GetSyncData gets the SyncData needed for syncing a remote gist and an associated local file.
func GetSyncData(gistFile conf.File, syncDataChan chan<- SyncData, config *conf.Conf) {
	fh, err := os.Open(gistFile.Path)
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken:  config.Gist.AccessToken,
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
			AccessToken:  config.Gist.AccessToken,
			Gist:         nil,
			GistFilename: "",
			FileContent:  nil,
			FilePath:     gistFile.Path,
			FileLastMod:  nil,
			Error:        err}
		return
	}
	fileLastMod := stat.ModTime()
	gist, resp, err := getGistData(config.Gist.AccessToken, gistFile.Id)
	if err != nil {
		syncDataChan <- SyncData{
			AccessToken:  config.Gist.AccessToken,
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
			AccessToken:  config.Gist.AccessToken,
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
			AccessToken:  config.Gist.AccessToken,
			Gist:         nil,
			GistFilename: "",
			FileContent:  nil,
			FilePath:     gistFile.Path,
			FileLastMod:  &fileLastMod,
			Error:        err}
		return
	}
	syncDataChan <- SyncData{
		AccessToken:  config.Gist.AccessToken,
		Gist:         gist,
		GistFilename: github.GistFilename(stat.Name()),
		FileContent:  fileContent,
		FilePath:     gistFile.Path,
		FileLastMod:  &fileLastMod,
		Error:        nil}
}

// Sync takes care of syncing the local file with the remote gist given
// SyncData.
func Sync(syncDataChan <-chan SyncData, syncResultChan chan<- SyncResult, config *conf.Conf) {
	data := <-syncDataChan
	gist := *data.Gist
	gistContent := gist.Files[data.GistFilename].Content
	if *gistContent == string(data.FileContent) {
		syncResultChan <- SyncResult{Error: nil, FileOverwritesGist: false, GistOverwritesFile: false}
		return
	}
	if data.FileLastMod.After(*gist.UpdatedAt) {
		_, _, err := updateGist(data.AccessToken, &gist, data.GistFilename, data.FileContent)
		if err != nil {
			syncResultChan <- SyncResult{Error: err, FileOverwritesGist: false, GistOverwritesFile: false}
			return
		}
		// TODO check for response.
		syncResultChan <- SyncResult{Error: nil, FileOverwritesGist: true, GistOverwritesFile: false}

	} else {
		gistContent := []byte(*gist.Files[data.GistFilename].Content)
		err := ioutil.WriteFile(data.FilePath, gistContent, 0644)
		if err != nil {
			syncResultChan <- SyncResult{Error: err, FileOverwritesGist: false, GistOverwritesFile: false}
			return
		}
		syncResultChan <- SyncResult{Error: nil, FileOverwritesGist: false, GistOverwritesFile: true}
	}
}
