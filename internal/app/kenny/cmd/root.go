package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/smf8/kenny/internal/app/kenny/cmd/devices"
	"github.com/smf8/kenny/internal/app/kenny/config"
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new kenny root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "kenny",
	}

	cfg := config.New()

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(cfg.Logger.Level)

	devices.Register(root)

	return root
}
