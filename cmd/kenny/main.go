package main

import (
	"os"

	"github.com/smf8/kenny/internal/app/kenny/cmd"
)

const exitCodeErr = 1

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitCodeErr)
		}
	}
}
