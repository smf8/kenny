package portaudio

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/gordonklaus/portaudio"
)

// POPlayer represents portaudio player. it implements audio.Player interface.
type POPlayer struct {
	streamParams portaudio.StreamParameters
	streams      []*audioStream
}

// NewPlayer creates a new POPlayer instance with given StreamSettings and deviceName.
// deviceName must be a valid portaudio device name from devices command.
// use DeviceNameDefault as deviceName to use system's default device.
func NewPlayer(deviceName string, settings StreamSettings) (*POPlayer, error) {
	param, err := newStreamParam(deviceName, PlayDeviceType, settings)
	if err != nil {
		return nil, err
	}

	return &POPlayer{
		streamParams: *param,
	}, nil
}

// OpenStream opens a new stream on the device and assigns an ID to it.
// you can use this ID to play, pause or stop the stream.
// you can send data to the channel (first dimension is for each channel) after calling POPlayer.Play()
//
// Make sure to close the stream after you're done using it.
func (p *POPlayer) OpenStream() (int, error) {
	buffer := make([]int16, p.streamParams.Output.Channels*p.streamParams.FramesPerBuffer)

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
		return -1, fmt.Errorf("failed to start playback %d: %w", id, err)
	}

	return id, nil
}

// CloseStream closes the stream and it's data channel.
func (p *POPlayer) CloseStream(streamID int) error {
	if err := p.streams[streamID].stream.Close(); err != nil {
		return fmt.Errorf("failed to close stream %d: %w", streamID, err)
	}

	return nil
}

// PausePlay pauses The play making next call to stream.Write() return an error.
func (p *POPlayer) PausePlay(streamID int) error {
	if err := p.streams[streamID].stream.Stop(); err != nil {
		return fmt.Errorf("failed to pause play %d: %w", streamID, err)
	}

	return nil
}

// ResumePlay starts stream processing in portaudio. you MUST call this between a PausePlay and Play call.
//
// calling this function for an already started stream will cause error.
func (p *POPlayer) ResumePlay(streamID int) error {
	if err := p.streams[streamID].stream.Start(); err != nil {
		return fmt.Errorf("failed to start play %d: %w", streamID, err)
	}

	return nil
}

// Play will add given data chunk to the output buffer. It will write the buffer to the output device once
// it is completely filled with data. use maximum chunk of `FramesPerBuffer`.
// data is copied into local buffer so using the same slice in multiple calls is safe.
func (p *POPlayer) Play(streamID int, data []int16) error {
	playStream := p.streams[streamID]

	if len(data) > len(playStream.buffer) {
		return fmt.Errorf("audio data size is larger than audio buffer. maximum allowed size is: %d", len(playStream.buffer))
	}

	if playStream.bufferIndex+len(data) >= len(playStream.buffer) {
		overflow := len(playStream.buffer) - playStream.bufferIndex
		copy(playStream.buffer[playStream.bufferIndex:], data)

		if err := playStream.stream.Write(); err != nil {
			logrus.Errorf("failed to write to output stream %d: %s", streamID, err.Error())
		}

		playStream.bufferIndex = copy(playStream.buffer, data[overflow:])
	} else {
		copy(playStream.buffer[playStream.bufferIndex:], data)
		playStream.bufferIndex += len(data)
	}

	return nil
}
