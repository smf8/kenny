package cmd

import (
	"github.com/spf13/cobra"
)

func main() error {
	return nil
}

// NewDevicesCommand creates a cobra command which will print list of recording audio devices
func NewDevicesCommand() *cobra.Command {
	devicesCmd := &cobra.Command{
		Use:   "kenny devices",
		Short: "this command will list available audio devices",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return main()
		},
	}

	return devicesCmd
}
