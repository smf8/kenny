package record

import (
	"bytes"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/smf8/kenny/internal/app/kenny/config"

	"github.com/smf8/kenny/internal/pkg/audio"
	"github.com/spf13/cobra"
)

//nolint:funlen,gomnd
func main(cfg config.Config) error {
	if err := audio.InitPortAudio(); err != nil {
		return err
	}

	defer audio.TerminatePortAudio()

	input, err := audio.DefaultInputDevice()
	if err != nil {
		return err
	}

	output, err := audio.DefaultOutputDevice()
	if err != nil {
		return err
	}

	storage := new(bytes.Buffer)
	settings := audio.StreamSettings{
		FramesPerBuffer:  cfg.Recorder.FramesPerBuffer,
		SampleRate:       cfg.Recorder.SampleRate,
		NumberOfChannels: cfg.Recorder.NumberOfChannels,
		Buffer:           storage,
	}

	inputStreamParam, err := audio.NewStreamParam(&settings, input, audio.RecordDeviceType)
	if err != nil {
		return err
	}

	outputSteamParam, err := audio.NewStreamParam(&settings, output, audio.PlayDeviceType)
	if err != nil {
		return err
	}

	api := audio.API{
		Settings: &settings,
	}

	err = api.OpenStream(inputStreamParam, audio.RecordDeviceType, api.SaveAudioCallback)
	if err != nil {
		return err
	}

	err = api.OpenStream(outputSteamParam, audio.PlayDeviceType, api.PlayAudioCallback)
	if err != nil {
		return err
	}

	go func() {
		if err := api.Record(); err != nil {
			log.Error(err)
		}
	}()

	<-time.After(100 * time.Millisecond)

	go func() {
		if err := api.Play(); err != nil {
			log.Error(err)
		}
	}()

	<-time.After(8 * time.Second)

	if err := api.PauseRecord(); err != nil {
		return err
	}

	if err := api.PausePlay(); err != nil {
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
