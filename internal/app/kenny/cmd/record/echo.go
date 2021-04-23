package record

import (
	"bytes"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/smf8/kenny/internal/app/kenny/audio"
	"github.com/smf8/kenny/internal/app/kenny/config"
	"github.com/smf8/kenny/internal/app/kenny/encoding"
)

//nolint:funlen,gomnd
func echo(cfg config.Config) error {
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

	encoder, err := encoding.NewEncoder(int(inputStreamParam.SampleRate),
		inputStreamParam.Input.Channels,
		cfg.Recorder.OpusFrameSizeMs)
	if err != nil {
		return err
	}

	decoder, err := encoding.NewDecoder(int(inputStreamParam.SampleRate),
		outputSteamParam.Output.Channels,
		cfg.Recorder.OpusFrameSizeMs)
	if err != nil {
		return err
	}

	err = api.OpenStream(inputStreamParam, audio.RecordDeviceType, encoder.CallbackRecord)
	if err != nil {
		return err
	}

	err = api.OpenStream(outputSteamParam, audio.PlayDeviceType, decoder.CallbackPlay)
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
		for {
			data, _ := encoder.Buffer.Read()

			_ = decoder.DecodeFrame(data)
		}
	}()
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
