package record

import (
	"bytes"

	"github.com/smf8/kenny/internal/app/kenny/config"

	"github.com/smf8/kenny/internal/pkg/audio"
	"github.com/spf13/cobra"
)

func main(cfg config.Config) error {
	if err := audio.InitPortAudio(); err != nil {
		return err
	}

	defer audio.TerminatePortAudio()

	input, err := audio.DefaultInputDevice()
	if err != nil {
		return err
	}

	buffer := make([]int32, 64)
	storage := new(bytes.Buffer)
	recorder := audio.Recorder{
		Buffer:           buffer,
		SampleRate:       cfg.Recorder.SampleRate,
		NumberOfChannels: cfg.Recorder.NumberOfChannels,
		InputDevice:      input,
		Storage:          storage,
	}

	if err := recorder.Record(); err != nil {
		return err
	}

	return nil
}

// Register registers record command to the root kenny command
//nolint:gomnd
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "record echo",
			Short: "this command will record something, encode it with opus, then decodes it and plays it back",
			Args:  cobra.ExactArgs(1),
			ValidArgs: []string{
				"echo",
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				return main(cfg)
			},
		},
	)
}
