package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/gordonklaus/portaudio"
	log "github.com/sirupsen/logrus"
)

// DeviceType represents audio device type. it's either RecordDeviceType or PlayDeviceType.
type DeviceType int

const (
	// RecordDeviceType is for audio input device.
	RecordDeviceType DeviceType = iota
	// PlayDeviceType is for audio output device.
	PlayDeviceType
)

// StreamSettings represent audio stream settings used for opening a portaudio stream
type StreamSettings struct {
	SampleRate       float64
	NumberOfChannels int
	Buffer           *bytes.Buffer
	FramesPerBuffer  int
}

// API represents audio stream api.
type API struct {
	input, output *portaudio.Stream
	Settings      *StreamSettings
	donePlay      chan interface{}
	doneRecord    chan interface{}
}

// InitPortAudio initializes portaudio and return an error if it failed.
func InitPortAudio() error {
	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize portaudio: %w", err)
	}

	log.Debugf(portaudio.VersionText())

	return nil
}

// TerminatePortAudio terminates portaudio instance. it should be called after we are finished working with pa API.
func TerminatePortAudio() {
	if err := portaudio.Terminate(); err != nil {
		log.Errorf("failed to terminate portaudio: %s", err.Error())
	}
}

// SaveAudioCallback is the callback function for saving audio to buffer
func (a *API) SaveAudioCallback(input []int32) {
	if err := binary.Write(a.Settings.Buffer, binary.BigEndian, input); err != nil {
		log.Errorf("failed to write byte to storage buffer: %s", err.Error())
	}
}

// PlayAudioCallback is the callback function to write audio buffer to stream output.
func (a *API) PlayAudioCallback(output []int32) {
	if err := binary.Read(a.Settings.Buffer, binary.BigEndian, output); err != nil {
		if errors.Is(err, io.EOF) {
			if err := a.PausePlay(); err != nil {
				log.Errorf("failed to pause audio play from callback: %s", err.Error())
			}

			return
		}

		log.Errorf("failed to read byte from buffer: %s", err.Error())
	}
}

// DefaultStreamParam returns default streamParameters which are being used by the host.
func DefaultStreamParam(deviceType DeviceType) (*portaudio.StreamParameters, error) {
	defaultHostAPI, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, fmt.Errorf("failed to get default host audio api: %w", err)
	}

	var param portaudio.StreamParameters

	switch deviceType {
	case RecordDeviceType:
		param = portaudio.LowLatencyParameters(
			defaultHostAPI.DefaultInputDevice,
			nil,
		)
	case PlayDeviceType:
		param = portaudio.LowLatencyParameters(
			nil,
			defaultHostAPI.DefaultOutputDevice,
		)
	}

	// for sake of simplicity, we only process mono audio
	param.Input.Channels = 1
	param.Input.Channels = 1
	param.FramesPerBuffer = 64

	return &param, nil
}

// NewStreamParam creates a new output/input stream which is ready to use.
// For temporarily pausing record/play use Stop() and to finish using stream call Close().
func NewStreamParam(s *StreamSettings, device *portaudio.DeviceInfo,
	deviceType DeviceType) (portaudio.StreamParameters, error) {
	deviceParameters := portaudio.StreamDeviceParameters{
		Device:   device,
		Channels: s.NumberOfChannels,
		Latency:  device.DefaultLowInputLatency,
	}
	sp := portaudio.StreamParameters{
		SampleRate:      s.SampleRate,
		FramesPerBuffer: s.FramesPerBuffer,
		Flags:           portaudio.ClipOff,
	}

	if deviceType == RecordDeviceType {
		sp.Input = deviceParameters
	} else if deviceType == PlayDeviceType {
		sp.Output = deviceParameters
	}

	return sp, nil
}

// OpenStream opens a stream with given callbackFunc for either input/output.
//for callbackFunc signature refer to portaudio.StreamCallback.
func (a *API) OpenStream(sp portaudio.StreamParameters,
	deviceType DeviceType, callbackFunc interface{}) error {
	stream, err := portaudio.OpenStream(sp, callbackFunc)
	if err != nil {
		return fmt.Errorf("failed to open stream for device type %d: %w", deviceType, err)
	}

	switch deviceType {
	case RecordDeviceType:
		a.input = stream
	case PlayDeviceType:
		a.output = stream
	}

	return nil
}

// Play starts the api's output stream. refer to stream's callbackFunc to handle audio.
// it's a blocking function.
func (a *API) Play() error {
	if a.donePlay == nil {
		a.donePlay = make(chan interface{})
	} else {
		return fmt.Errorf("an instance is already using Play")
	}

	if err := a.output.Start(); err != nil {
		return fmt.Errorf("failed to start play stream: %w", err)
	}

	log.Debugf("playing...")
	<-a.donePlay

	return nil
}

// PausePlay stops the output stream.
func (a *API) PausePlay() error {
	if err := a.output.Stop(); err != nil {
		return fmt.Errorf("failed to stop play stream: %w", err)
	}

	close(a.donePlay)

	log.Debugf("pause play")

	return nil
}

// FinishPlay closes API's output stream.
func (a *API) FinishPlay() error {
	if err := a.output.Close(); err != nil {
		return fmt.Errorf("failed to close play stream: %w", err)
	}

	a.donePlay = nil

	return nil
}

// Record starts the api's input stream. refer to stream's callbackFunc to handle audio.
// it's a blocking function.
func (a *API) Record() error {
	if a.doneRecord == nil {
		a.doneRecord = make(chan interface{})
	}

	if err := a.input.Start(); err != nil {
		return fmt.Errorf("failed to start record stream: %w", err)
	}

	log.Debugf("recording...")
	<-a.doneRecord

	return nil
}

//PauseRecord stops the input stream.
func (a *API) PauseRecord() error {
	if err := a.input.Stop(); err != nil {
		return fmt.Errorf("failed to stop record stream: %w", err)
	}

	log.Debugf("pause recording")
	close(a.doneRecord)

	return nil
}

// FinishRecord closes API's input stream.
func (a *API) FinishRecord() error {
	if err := a.input.Close(); err != nil {
		return fmt.Errorf("failed to close play stream: %w", err)
	}

	return nil
}
