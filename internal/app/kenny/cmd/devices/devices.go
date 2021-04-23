package devices

import (
	"fmt"

	"github.com/smf8/kenny/internal/app/kenny/audio"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(deviceArg string) error {
	if err := audio.InitPortAudio(); err != nil {
		return err
	}

	defer audio.TerminatePortAudio()

	switch deviceArg {
	case "list":
		devices, err := audio.DevicesInfo()
		if err != nil {
			log.Errorf("failed to get audio device list: %s", err.Error())
		}

		for i := range devices {
			fmt.Println(devices[i])
		}
	case "input":
		device, err := audio.DefaultInputDevice()
		if err != nil {
			log.Errorf("failed to get default input audio device: %s", err.Error())
		}

		fmt.Println(audio.PrintDeviceInfo(device))
	case "output":
		device, err := audio.DefaultOutputDevice()
		if err != nil {
			log.Errorf("failed to get default output audio device: %s", err.Error())
		}

		fmt.Println(audio.PrintDeviceInfo(device))
	}

	return nil
}

// Register registers devices command to the root kenny command
//nolint:gomnd
func Register(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:   "devices {list | input | output}",
			Short: "this command will list available audio devices",
			Args:  cobra.ExactArgs(1),
			ValidArgs: []string{
				"list",
				"input",
				"output",
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				return main(args[0])
			},
		},
	)
}
