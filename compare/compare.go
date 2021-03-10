package compare

import (
	"log"
	"os"

	"github.com/kahunacohen/sinker/conf"

	"github.com/kahunacohen/sinker/gist"
)

func Compare(config conf.Conf, file conf.File, which chan *gist.SyncData) {
	fh, err := os.Open(file.Path)
	if err != nil {
		log.Fatalf("problem reading file: %s", err)
	}
	resp := gist.GetSyncData(config.Gist.AccessToken, fh, file.Id)
	which <- &resp
}
