package gist

import (
	"testing"
)

func TestGet(t *testing.T) {
	// Access token, gist ID
	gist, resp, err := Get("8a90f5cd571f737efa8fe5cdb677e5ab86c2e5ff", "142a4dfb66f0e2eab38cb68e0b69d95c")
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
