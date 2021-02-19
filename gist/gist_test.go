package gist

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Intercept http request to github.
	httpmock.RegisterResponder("GET", "https://api.github.com/gists/142a4dfb66f0e2eab38cb68e0b69d95c",
		httpmock.NewStringResponder(200, `{"files": {".bashrc": {"filename": ".bashrc"}}}`))

	// Pass Access token and gist ID
	gist, resp, err :=
		Get(os.Getenv("SINKER_GIST_ACCESS_TOKEN"), "142a4dfb66f0e2eab38cb68e0b69d95c")
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
