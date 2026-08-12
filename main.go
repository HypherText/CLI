package main

import (
	"flag"
	"os"
	"path/filepath"
)

func main() {
	builder := Builder{}
	flag.StringVar(&builder.workdir, "pwd", "", "Set the working directory, defaults to the current working directory.")
	flag.Parse()

	if builder.workdir == "" {
		builder.workdir, _ = os.Getwd()
	}
	builder.workdir, _ = filepath.Abs(builder.workdir)
	builder.BuildAssets()
}
