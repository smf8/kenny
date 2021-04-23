package encoding

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/hraban/opus.v2"
)

const (
	bufferSize      = 1024
	frameBufferSize = 1000
)

// Encoder represents an opus audio encoder.
type Encoder struct {
	Buffer     *FrameBuffer
	E          *opus.Encoder
	SampleRate int
	Channels   int

	// opus encoder only accepts PCM audio data of size 2.5, 5, 10, 20, 40, 60 ms
	// but portaudio callback function is called whenever it has captured audio data
	// so we use a temporary buffer to store one frame.
	pcmBuffer []int16
	pcmCursor int
	FrameSize int
}

// NewEncoder creates a new encoder instance with a new buffer. You can access encoded data from buffer.
func NewEncoder(sampleRate, channels, frameSize int) (*Encoder, error) {
	e, err := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)
	if err != nil {
		return nil, fmt.Errorf("failed to create opus encoder: %w", err)
	}

	return &Encoder{
		Buffer:     NewFrameBuffer(frameBufferSize),
		E:          e,
		SampleRate: sampleRate,
		Channels:   channels,
		pcmBuffer:  make([]int16, frameSize),
	}, nil
}

// CallbackRecord handles recording and encoding procedures.
// input pcm data is recorded and when it reaches Encoder.FrameSize
// it will be written into the buffer as an opus frame
func (e *Encoder) CallbackRecord(input []int16) {
	if e.pcmCursor+len(input) >= len(e.pcmBuffer) {
		// pcmBuffer is full, encode a frame
		overflow := len(e.pcmBuffer) - e.pcmCursor
		copy(e.pcmBuffer[e.pcmCursor:], input[:overflow])
		e.pcmCursor = 0

		// encode e.pcmBuffer
		b := make([]byte, bufferSize)
		n, err := e.E.Encode(e.pcmBuffer, b)

		if err != nil {
			log.Errorf("failed to encode pcm data: %s", err.Error())
		}

		e.Buffer.Write(b[:n])
	} else {
		copy(e.pcmBuffer[e.pcmCursor:], input)
		e.pcmCursor += len(input)
	}
}
