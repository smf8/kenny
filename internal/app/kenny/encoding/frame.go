package encoding

import (
	"io"
)

// FrameBuffer represents opus data frame. each frame contains a variable sized []byte data,
// and providing slices of data with different size leads to decoding error.
// so a circular array of []byte are used to save opus audio frames.
type FrameBuffer struct {
	buffer                [][]byte
	readIndex, writeIndex int
	size                  int
}

//NewFrameBuffer creates a new FrameBuffer.
func NewFrameBuffer(size int) *FrameBuffer {
	b := make([][]byte, size)

	return &FrameBuffer{
		buffer: b,
		size:   size,
	}
}

// Write writes data to buffer at writeIndex.
func (f *FrameBuffer) Write(data []byte) {
	f.buffer[f.writeIndex] = make([]byte, len(data))
	copy(f.buffer[f.writeIndex], data)
	f.writeIndex = (f.writeIndex + 1) % f.size
}

// Read returns data currently located in readIndex cursor.
func (f *FrameBuffer) Read() ([]byte, error) {
	data := f.buffer[f.readIndex]
	if data == nil {
		return nil, io.EOF
	}

	f.readIndex = (f.readIndex + 1) % f.size

	return data, nil
}
