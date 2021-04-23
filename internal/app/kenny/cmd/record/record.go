package record

import (
	"github.com/smf8/kenny/internal/app/kenny/config"

	"github.com/spf13/cobra"
)

const (
	commandEcho   = "echo"
	encodeCommand = "encode"
)

// Register registers record command to the root kenny command
//nolint:gomnd
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "record {echo | encode}",
			Short: "this command will record something, encode it with opus, then decodes it and plays it back",
			Args:  cobra.ExactArgs(1),
			ValidArgs: []string{
				"echo",
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				switch args[0] {
				case commandEcho:
					return echo(cfg)
				}

				return nil
			},
		},
	)
}
