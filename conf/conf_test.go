package conf

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseJsonConfig(t *testing.T) {
	jsonByte := []byte(`{
		"gist": {
			"accessToken": "xxx",
			"files": [
				{
					"path": "/a/b",
					 "id": "123"
				}
			]
		}
	}`)
	config, err := ParseJsonConfg(jsonByte)
	if err != nil {
		t.Fatalf("Didn't expect error, but got %s", err)
	}
	if config.Gist.AccessToken != "xxx" {
		t.Fatalf("Expected access token to equal xxx but got %v", config.Gist.AccessToken)
	}
	got := config.Gist.Files
	lenGot := len(got)
	if lenGot != 1 {
		t.Fatalf("Expected files to be length 2, but got %d", lenGot)
	}
	file := got[0]
	if file.Path != "/a/b" || file.Id != "123" {
		t.Errorf("the file is not what we expect. Got: Path: %s and Id: %s", file.Path, file.Id)
	}
}

func TestParseJsonConfigBadParse(t *testing.T) {
	jsonByte := []byte(`{
		"gist: {
			"accessToken": "xxx",
			"files": []
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
func TestGet(t *testing.T) {
	sinkerRcPath := "./sinkerrc.json"
	jsonByte := []byte(`{
		"gist": {
			"accessToken": "xxx",
			"files": []
		}
	}`)
	err := ioutil.WriteFile(sinkerRcPath, jsonByte, 0664)
	defer os.Remove(sinkerRcPath)
	if err != nil {
		t.Fatalf("not able to create test config file: %s", err)
	}

	conf, err := Get(sinkerRcPath)
	if err != nil {
		t.Errorf("problem getting or parsing config file: %s", err)
	}
	want := "xxx"
	got := conf.Gist.AccessToken
	if got != want {
		t.Errorf("want access token %s, got %s", want, got)
	}

}
