package portaudio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// PORecorder represents portaudio recorder. it implements audio.Recorder interface.
type PORecorder struct {
	streamParams portaudio.StreamParameters
	streams      []recordStream
}

// because we may have multiple audio sources(i.e multiple speakers),
// we represent each one with a recordStream instance.
type recordStream struct {
	stream  *portaudio.Stream
	channel chan [][]int16
	buffer  [][]int16
}

// NewRecorder creates a new PORecorder instance with given StreamSettings.
// deviceName must be a valid portaudio device name from devices command.
func NewRecorder(deviceName string, settings StreamSettings) (*PORecorder, error) {
	param, err := newStreamParam(deviceName, RecordDeviceType, settings)
	if err != nil {
		return nil, err
	}

	buffer := make([][]int16, settings.Channels)
	for i := range buffer {
		buffer[i] = make([]int16, settings.FramesPerBuffer)
	}

	return &PORecorder{
		streamParams: *param,
	}, nil
}

// OpenStream opens a new stream on the device and assigns an ID to it.
// you can use this ID to record, pause or stop the stream.
// you can receive audio data (first dimension is the number of audio channels) from
// the channel after you call Record.
//
// Make sure to close the stream after you're done using it.
func (p *PORecorder) OpenStream() (int, <-chan [][]int16, error) {
	buffer := make([][]int16, p.streamParams.Input.Channels)
	for i := range buffer {
		buffer[i] = make([]int16, p.streamParams.FramesPerBuffer)
	}

	stream, err := portaudio.OpenStream(p.streamParams, buffer)
	if err != nil {
		return -1, nil, fmt.Errorf("failed to open stream: %w", err)
	}

	if err := stream.Start(); err != nil {
		return -1, nil, fmt.Errorf("failed to start stream: %w", err)
	}

	ch := make(chan [][]int16)
	id := len(p.streams)
	s := recordStream{
		stream:  stream,
		channel: ch,
		buffer:  buffer,
	}
	p.streams = append(p.streams, s)

	return id, ch, nil
}

// CloseStream closes the stream and it's data channel.
func (p *PORecorder) CloseStream(streamID int) error {
	if err := p.streams[streamID].stream.Close(); err != nil {
		return fmt.Errorf("failed to close stream %d: %w", streamID, err)
	}

	close(p.streams[streamID].channel)

	return nil
}

// PauseRecord pauses The record making next call to stream.Read() return an error.
func (p *PORecorder) PauseRecord(streamID int) error {
	if err := p.streams[streamID].stream.Stop(); err != nil {
		return fmt.Errorf("failed to pause record %d: %w", streamID, err)
	}

	return nil
}

// Record starts audio recording inside a go routine.
// sending captured data into stream's channel.
func (p *PORecorder) Record(streamID int) error {
	recordStream := p.streams[streamID]
	if err := recordStream.stream.Start(); err != nil {
		return fmt.Errorf("failed to start record %d: %w", streamID, err)
	}

	var err error

	go func() {
		for {
			if err = recordStream.stream.Read(); err != nil {
				err = fmt.Errorf("failed to read input stream %d: %w", streamID, err)
				break
			}

			recordStream.channel <- recordStream.buffer
		}
	}()

	return nil
}
