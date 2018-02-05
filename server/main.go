package main

import (
	"github.com/ArthurHlt/gridana"
	_ "github.com/ArthurHlt/gridana/drivers/alertmanager"
	"os"
	"time"
)

var version string

func main() {
	if version == "" {
		version = time.Now().Format(time.RFC3339) + "-build"
	}
	gridana.NewAppWithVersion(version).Run(os.Args)
}
