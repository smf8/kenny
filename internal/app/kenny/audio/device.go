package audio

import (
	"fmt"
	"strings"

	"github.com/gordonklaus/portaudio"
)

// DevicesInfo returns the list of available audio devices.
// `portaudio.Initialize()` must be called before running this function
func DevicesInfo() ([]string, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	for _, d := range devices {
		result = append(result, PrintDeviceInfo(d))
	}

	return result, nil
}

// DefaultInputDevice returns the default recording device.
// `portaudio.Initialize()` must be called before running this function
func DefaultInputDevice() (*portaudio.DeviceInfo, error) {
	device, err := portaudio.DefaultInputDevice()
	if err != nil {
		return nil, err
	}

	return device, nil
}

// DefaultOutputDevice returns the default recording device.
// `portaudio.Initialize()` must be called before running this function.
func DefaultOutputDevice() (*portaudio.DeviceInfo, error) {
	device, err := portaudio.DefaultOutputDevice()
	if err != nil {
		return nil, err
	}

	return device, nil
}

// PrintDeviceInfo formats and prints portaudio.DeviceInfo.
func PrintDeviceInfo(info *portaudio.DeviceInfo) string {
	sb := new(strings.Builder)

	fmt.Fprintf(sb, "==========================================\n")
	fmt.Fprintf(sb, "[Name]: %s\n", info.Name)
	fmt.Fprintf(sb, "[Max input Channels]: %d\n", info.MaxInputChannels)
	fmt.Fprintf(sb, "[Max output Channels]: %d\n", info.MaxOutputChannels)
	fmt.Fprintf(sb, "[Default Low input Latency]: %s\n", info.DefaultLowInputLatency)
	fmt.Fprintf(sb, "[Default High input Latency]: %s\n", info.DefaultHighInputLatency)
	fmt.Fprintf(sb, "[Default Low output Latency]: %s\n", info.DefaultLowOutputLatency)
	fmt.Fprintf(sb, "[Default High output Latency]: %s\n", info.DefaultHighOutputLatency)
	fmt.Fprintf(sb, "[Default Sample Rate]: %f\n", info.DefaultSampleRate)

	if info.HostApi != nil {
		fmt.Fprintf(sb, "[HOST API Info]: \n")
		fmt.Fprintf(sb, "\t[Type]: %s\n", info.HostApi.Type)
		fmt.Fprintf(sb, "\t[Name]: %s\n", info.HostApi.Name)

		if info.HostApi.DefaultInputDevice != nil {
			fmt.Fprintf(sb, "\t[Default InputDevice Name]: %s\n", info.HostApi.DefaultInputDevice.Name)
		}

		if info.HostApi.DefaultOutputDevice != nil {
			fmt.Fprintf(sb, "\t[Default OutputDevice Name]: %s\n", info.HostApi.DefaultOutputDevice.Name)
		}
	}

	return sb.String()
}
