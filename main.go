package main

import (
	"github.com/spf13/afero"
	"os"
	"vstruct/cmd"
)

func main() {
	cmd.Execute(afero.NewOsFs(), os.Args)
}
