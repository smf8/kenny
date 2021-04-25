package audio

// API represents the wrapper for playback and recording.
type API struct {
	Player   Player
	Recorder Recorder
}

// Player is a general interface for audio playback.
// call OpenStream for each media stream and close it
// after you're done using it.
type Player interface {
	OpenStream() (int, chan<- [][]int16, error)
	Play(streamID int) error
	PausePlay(streamID int) error
	CloseStream(streamID int) error
}

// Recorder is a general interface for audio recording.
// call OpenStream for each media input (i.e. microphone and system sound)
// and close it after you're done using it.
type Recorder interface {
	OpenStream() (int, <-chan [][]int16, error)
	Record(streamID int) error
	PauseRecord(streamID int) error
	CloseStream(streamID int) error
}
