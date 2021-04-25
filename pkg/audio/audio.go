package audio

type API struct {
	Player   Player
	Recorder Recorder
}

type Player interface {
	OpenStream() (int, chan<- [][]int16, error)
	Play(streamID int) error
	PausePlay(streamID int) error
	CloseStream(streamID int) error
}

type Recorder interface {
	OpenStream() (int, <-chan [][]int16, error)
	Record(streamID int) error
	PauseRecord(streamID int) error
	CloseStream(streamID int) error
}
