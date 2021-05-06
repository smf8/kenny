package portaudio

import (
	"errors"
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// DeviceType represents audio device type. it's either RecordDeviceType or PlayDeviceType.
type DeviceType int

const (
	// RecordDeviceType is for audio input device.
	RecordDeviceType DeviceType = iota
	// PlayDeviceType is for audio output device.
	PlayDeviceType

	// DeviceNameDefault is the default devicename which lets port audio api to choose default audio device.
	DeviceNameDefault = "default"
)

var (
	// ErrDeviceNotFound occurs when device with given info does not exist in the system.
	ErrDeviceNotFound = errors.New("failed to find audio device")
)

// StreamSettings represent audio stream settings used for opening a portaudio stream
type StreamSettings struct {
	SampleRate      float64
	Channels        int
	FramesPerBuffer int
}

// because we may have multiple audio sources(i.e multiple speakers),
// we represent each one with a audioStream instance.
type audioStream struct {
	stream      *portaudio.Stream
	buffer      []int16
	bufferIndex int
}

// Init must be called before any further usage on portaudio package
func Init() error {
	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize portaudio: %w", err)
	}

	return nil
}

// Cleanup must be called when you are done using portaudio package.
func Cleanup() error {
	if err := portaudio.Terminate(); err != nil {
		return fmt.Errorf("failed to terminate portaudio: %w", err)
	}

	return nil
}

func newStreamParam(
	deviceName string, deviceType DeviceType, settings StreamSettings,
) (*portaudio.StreamParameters, error) {
	defaultHostAPI, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, fmt.Errorf("failed to get default host api: %w", err)
	}

	var device *portaudio.DeviceInfo

	for _, d := range defaultHostAPI.Devices {
		if d.Name == deviceName {
			device = d
		}
	}

	if deviceName == DeviceNameDefault {
		if deviceType == RecordDeviceType {
			device = defaultHostAPI.DefaultInputDevice
		} else {
			device = defaultHostAPI.DefaultOutputDevice
		}
	}

	if device == nil {
		return nil, ErrDeviceNotFound
	}

	deviceParameters := portaudio.StreamDeviceParameters{
		Device:   device,
		Channels: settings.Channels,
	}
	sp := portaudio.StreamParameters{
		SampleRate:      settings.SampleRate,
		FramesPerBuffer: settings.FramesPerBuffer,
		Flags:           portaudio.ClipOff,
	}

	switch deviceType {
	case RecordDeviceType:
		deviceParameters.Latency = device.DefaultLowInputLatency
		sp.Input = deviceParameters
	case PlayDeviceType:
		deviceParameters.Latency = device.DefaultLowOutputLatency
		sp.Output = deviceParameters
	}

	return &sp, nil
}

// NewStreamParam creates a new output/input stream which is ready to use.
// For temporarily pausing record/play use Stop() and to finish using stream call Close().
func NewStreamParam(s *StreamSettings, device *portaudio.DeviceInfo,
	deviceType DeviceType) (portaudio.StreamParameters, error) {
	deviceParameters := portaudio.StreamDeviceParameters{
		Device:   device,
		Channels: s.Channels,
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

func openStream(buffer []int16, params portaudio.StreamParameters) (*portaudio.Stream, error) {
	stream, err := portaudio.OpenStream(params, buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}

	return stream, nil
}
