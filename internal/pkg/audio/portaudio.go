package audio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
	log "github.com/sirupsen/logrus"
)

// InitPortAudio initializes portaudio and return an error if it failed.
func InitPortAudio() error {
	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize portaudio: %w", err)
	}

	log.Debugf(portaudio.VersionText())

	return nil
}

// TerminatePortAudio terminates portaudio instance. it should be called after we are finished working with pa API.
func TerminatePortAudio() {
	if err := portaudio.Terminate(); err != nil {
		log.Errorf("failed to terminate portaudio: %s", err.Error())
	}
}
