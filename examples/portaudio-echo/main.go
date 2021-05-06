package main

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/smf8/kenny/pkg/audio/portaudio"
)

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := portaudio.Init(); err != nil {
		panic(err)
	}

	defer portaudio.Cleanup()

	settings := portaudio.StreamSettings{
		SampleRate:      48000,
		Channels:        2,
		FramesPerBuffer: 500,
	}

	recorder, err := portaudio.NewRecorder(portaudio.DeviceNameDefault, settings)
	if err != nil {
		panic(err)
	}

	recordID, err := recorder.OpenStream()
	if err != nil {
		panic(err)
	}

	defer recorder.CloseStream(recordID)

	player, err := portaudio.NewPlayer(portaudio.DeviceNameDefault, settings)
	if err != nil {
		panic(err)
	}

	playID, err := player.OpenStream()
	if err != nil {
		panic(err)
	}

	defer player.CloseStream(playID)

	// 1 * 2(channels) * sample rate = 1s delay.
	latencyBuffer := make([]int16, int(settings.SampleRate)*settings.Channels)
	index := 0
	chunk := make([]int16, settings.FramesPerBuffer*settings.Channels)

	go func() {
		for {

			a, err := recorder.Record(recordID)
			if err != nil {
				log.Errorln(err)
			}

			for i := range a {

				chunk[i] = latencyBuffer[index]
				latencyBuffer[index] = a[i]

				index = (index + 1) % len(latencyBuffer)
			}

			if err := player.Play(playID, chunk[:len(a)]); err != nil {
				log.Errorln(err)
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	<-sig
}
