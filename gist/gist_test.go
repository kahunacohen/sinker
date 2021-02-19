package gist

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Intercept http request to github.
	httpmock.RegisterResponder("GET", "https://api.github.com/gists/142a4dfb66f0e2eab38cb68e0b69d95c",
		httpmock.NewStringResponder(200, `{"files": {".bashrc": {"filename": ".bashrc"}}}`))

	// Pass Access token (doesn't matter what it is since we are mocking http request) and gist ID
	gist, resp, err :=
		Get("xxx", "142a4dfb66f0e2eab38cb68e0b69d95c")
	if err != nil {
		t.Fatalf("wanted, no error, got: %s", err)
	}
	got := *(gist.Files[".bashrc"].Filename)
	if got != ".bashrc" {
		t.Errorf("wanted .bashrc, got %s", err)
	}
	gotStatus := (*resp.Response).StatusCode

	if gotStatus != 200 {
		t.Errorf("wanted 200 status code, got: %s", got)
	}
}
