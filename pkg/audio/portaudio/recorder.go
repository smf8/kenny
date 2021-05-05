package portaudio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// PORecorder represents portaudio recorder. it implements audio.Recorder interface.
type PORecorder struct {
	streamParams portaudio.StreamParameters
	streams      []*audioStream
}

// NewRecorder creates a new PORecorder instance with given StreamSettings.
// deviceName must be a valid portaudio device name from devices command.
// use DeviceNameDefault as deviceName to use system's default device.
func NewRecorder(deviceName string, settings StreamSettings) (*PORecorder, error) {
	param, err := newStreamParam(deviceName, RecordDeviceType, settings)
	if err != nil {
		return nil, err
	}

	return &PORecorder{
		streamParams: *param,
	}, nil
}

// OpenStream opens a new stream on the device and assigns an ID to it.
// you can use this ID to record, pause or stop the stream.
// This function will also Start the portaudio stream so Make sure to close it after you are done using it.
func (p *PORecorder) OpenStream() (int, error) {
	buffer := make([]int16, p.streamParams.Input.Channels*p.streamParams.FramesPerBuffer)

	stream, err := openStream(buffer, p.streamParams)
	if err != nil {
		return -1, err
	}

	id := len(p.streams)
	s := &audioStream{
		stream: stream,
		buffer: buffer,
	}
	p.streams = append(p.streams, s)

	if err := stream.Start(); err != nil {
		return -1, fmt.Errorf("failed to start record %d: %w", id, err)
	}

	return id, nil
}

// CloseStream closes the stream and it's data channel.
func (p *PORecorder) CloseStream(streamID int) error {
	if err := p.streams[streamID].stream.Close(); err != nil {
		return fmt.Errorf("failed to close stream %d: %w", streamID, err)
	}

	return nil
}

// PauseRecord pauses The record making next call to stream.Read() return an error.
//
//
func (p *PORecorder) PauseRecord(streamID int) error {
	if err := p.streams[streamID].stream.Stop(); err != nil {
		return fmt.Errorf("failed to pause record %d: %w", streamID, err)
	}

	return nil
}

// ResumeRecord starts stream processing in portaudio. you MUST call this between a PauseRecord and Record call.
//
// calling this function for an already started stream will cause error.
func (p *PORecorder) ResumeRecord(streamID int) error {
	if err := p.streams[streamID].stream.Start(); err != nil {
		return fmt.Errorf("failed to start record %d: %w", streamID, err)
	}

	return nil
}

// Record reads an audio chunk from device and returns it.
//
// the returned slice will be changed in next calls to Record, So use copy to store it.
func (p *PORecorder) Record(streamID int) ([]int16, error) {
	recordStream := p.streams[streamID]

	if err := recordStream.stream.Read(); err != nil {
		return nil, fmt.Errorf("failed to read input stream %d: %w", streamID, err)
	}

	return recordStream.buffer, nil
}
