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
)

var (
	//ErrDeviceNotFound occurs when device with given info does not exist in the system.
	ErrDeviceNotFound = errors.New("failed to find audio device")
)

// StreamSettings represent audio stream settings used for opening a portaudio stream
type StreamSettings struct {
	SampleRate      float64
	Channels        int
	FramesPerBuffer int
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

	if device == nil {
		return nil, ErrDeviceNotFound
	}

	var param portaudio.StreamParameters

	switch deviceType {
	case RecordDeviceType:
		param = portaudio.LowLatencyParameters(device, nil)
		param.Input.Channels = settings.Channels
	case PlayDeviceType:
		param = portaudio.LowLatencyParameters(nil, device)
		param.Output.Channels = settings.Channels
	}

	param.FramesPerBuffer = settings.FramesPerBuffer
	param.SampleRate = settings.SampleRate

	return &param, nil
}
