package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"gopkg.in/hraban/opus.v2"
)

// Decoder is used to decode opus encoded data into PCM.
type Decoder struct {
	Buffer     *bytes.Buffer
	D          *opus.Decoder
	SampleRate int
	Channels   int
	FrameSize  int
}

// NewDecoder creates a Decoder instance.
// buffer is the encoded data buffer. call Decode to start decoding data Decoder.Buffer.
func NewDecoder(sampleRate, channels, frameSize int) (*Decoder, error) {
	d, err := opus.NewDecoder(sampleRate, channels)
	if err != nil {
		return nil, fmt.Errorf("failed to create opus decoder: %w", err)
	}

	return &Decoder{
		Buffer:     new(bytes.Buffer),
		D:          d,
		SampleRate: sampleRate,
		Channels:   channels,
		FrameSize:  frameSize,
	}, nil
}

// DecodeFrame decodes given frame to pcm data, the frame must be a valid opus frame data
// encoded with Decoder.FrameSize. decoded pcm data is saved in Decoder.Buffer
func (d *Decoder) DecodeFrame(frame []byte) error {
	pcmData := make([]int16, d.FrameSize)

	n, err := d.D.Decode(frame, pcmData)
	if err != nil {
		return fmt.Errorf("failed to decode opus data to pcm: %w", err)
	}

	err = binary.Write(d.Buffer, binary.BigEndian, pcmData[:n])
	if err != nil {
		return fmt.Errorf("failed to save pcm data to decoder buffer: %w", err)
	}

	return nil
}

//CallbackPlay is a portaudio callback function to decode and write Data to output device
func (d *Decoder) CallbackPlay(output []int16) {
	if err := binary.Read(d.Buffer, binary.BigEndian, output); err != nil {
		for i := range output {
			output[i] = 0
		}

		if err != io.EOF {
			log.Errorf("failed to play during callback: %s", err.Error())
		}
	}
}
