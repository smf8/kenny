package main

import (
	"github.com/smf8/kenny/internal/app/kenny/cmd"
	"os"
)

func main() {
	root := cmd.NewDevicesCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(1)
		}
	}
}
