package main

import (
	"github.com/ArthurHlt/gridana"
	_ "github.com/ArthurHlt/gridana/drivers/alertmanager"
	"os"
)

func main() {
	gridana.NewApp().Run(os.Args)
}
