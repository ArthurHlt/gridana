// +build !ui

package gridana

import (
	"net/http"
	"os"
)

type NoUI struct{}

func (NoUI) Open(name string) (http.File, error) {
	return nil, os.ErrNotExist
}

func assetFS() http.FileSystem {
	return NoUI{}
}
