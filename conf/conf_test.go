package conf

import (
	"testing"
)

func TestParseJsonConfig(t *testing.T) {
	jsonByte := []byte(`{
		"gist": {
			"accessToken": "xxx",
			"files": ["a", "b"]
		}
	}`)
	config, err := ParseJsonConfg(jsonByte)
	if err != nil {
		t.Fatalf("Didn't expect error, but got %s", err)
	}
	if config.Gist.AccessToken != "xxx" {
		t.Errorf("Expected access token to equal xxx but got %v", config.Gist.AccessToken)
	}
	got := config.Gist.Files
	lenGot := len(got)
	if lenGot != 2 {
		t.Errorf("Expected files to be length 2, but got %d", lenGot)
	}
	if got[0] != "a" || got[1] != "b" {
		t.Errorf("One of the files is not what we expect. Got: %s, %s", got[0], got[1])
	}
}

func TestParseJsonConfigBadParse(t *testing.T) {
	jsonByte := []byte(`{
		"gist: {
			"accessToken": "xxx",
			"files": ["a", "b"]
		}
	}`)
	config, err := ParseJsonConfg(jsonByte)
	if config.Gist.AccessToken != "" && config.Gist.Files != nil {
		t.Error("Access token and files are not zero value")
	}
	errorWant := `invalid character '\n' in string literal`
	errorGot := err.Error()
	if errorGot != errorWant {
		t.Errorf("wanted \"%s\" for error, got \"%s\"", errorWant, errorGot)
	}
}
