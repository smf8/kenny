package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"
	log "github.com/sirupsen/logrus"
)

// Recorder is used to record from an audioInput device. It also contains information
// for recorded data storage
type Recorder struct {
	InputDevice      *portaudio.DeviceInfo
	SampleRate       float64
	NumberOfChannels int
	Buffer           []int32
	Storage          *bytes.Buffer
	NSamples         int
}

func (r *Recorder) saveInputStream(input []int32) {
	if err := binary.Write(r.Storage, binary.BigEndian, input); err != nil {
		log.Errorf("failed to write byte to storage buffer: %s", err.Error())
	}

	r.NSamples += len(input)
}

// Record creates and starts an audio stream with a custom StreamCallback function.
func (r *Recorder) Record() error {
	inputParameters := portaudio.StreamDeviceParameters{
		Device:   r.InputDevice,
		Channels: r.NumberOfChannels,
		Latency:  r.InputDevice.DefaultLowInputLatency,
	}

	sp := portaudio.StreamParameters{
		Input:           inputParameters,
		SampleRate:      r.SampleRate,
		FramesPerBuffer: len(r.Buffer),
		Flags:           portaudio.ClipOff,
	}

	stream, err := portaudio.OpenStream(sp, r.saveInputStream)
	if err != nil {
		return fmt.Errorf("failed to open record stream: %w", err)
	}

	defer func() {
		if err := stream.Stop(); err != nil {
			log.Errorf("failed to stop record stream: %s", err.Error())
		} else if err = stream.Close(); err != nil {
			log.Errorf("failed to close record stream: %s", err.Error())
		}
	}()

	if err := stream.Start(); err != nil {
		return fmt.Errorf("failed to start record stream: %w", err)
	}

	log.Debugf("recording started...\n")
	//nolint:gomnd
	<-time.After(4 * time.Second)
	log.Debugf("finished recording.\n")

	r.playAudio()

	return nil
}

// TODO: remove this function
//nolint: gomnd
func (r *Recorder) playAudio() {
	out := make([]int32, 8192)

	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), out)
	if err != nil {
		log.Errorf("failed to open output stream: %s", err.Error())

		return
	}

	if err := stream.Start(); err != nil {
		log.Errorf("failed to start stream: %s", err.Error())

		return
	}

	defer func() {
		if err := stream.Stop(); err != nil {
			log.Errorf("failed to start stream: %s", err.Error())

			return
		}
	}()

	for remaining := r.NSamples; remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}

		if err := binary.Read(r.Storage, binary.BigEndian, out); err != nil {
			log.Errorf("failed to read byte from storage: %s", err.Error())
		}

		if err := stream.Write(); err != nil {
			log.Errorf("failed to write stream: %s", err.Error())
		}

		fmt.Println(remaining, len(out))
	}
}
