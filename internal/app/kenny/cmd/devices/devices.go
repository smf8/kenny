package devices

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
	log "github.com/sirupsen/logrus"
	"github.com/smf8/kenny/internal/pkg/audio"
	"github.com/spf13/cobra"
)

func main(deviceArg string) error {
	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize portaudio: %w", err)
	}

	defer func() {
		if err := portaudio.Terminate(); err != nil {
			log.Debugf("failed to terminate portaudio: %s", err.Error())
		}
	}()

	log.Debugf(portaudio.VersionText())

	switch deviceArg {
	case "list":
		devices, err := audio.AllDevices()
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

		fmt.Println(device)
	case "output":
		device, err := audio.DefaultOutputDevice()
		if err != nil {
			log.Errorf("failed to get default output audio device: %s", err.Error())
		}

		fmt.Println(device)
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
