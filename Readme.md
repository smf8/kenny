# Kenny
I'm just trying to make a cli operated voice call chat application using go with help of [webRTC](https://github.com/pion) and [PortAudio](https://github.com/gordonklaus/portaudio/).

It might stay a Work In Progress for a long time.

## Usage
`go get` the project or clone it and use `go mod download` to install dependencies. You also must [install portaudio](http://portaudio.com/docs/v19-doxydocs/tutorial_start.html).

after that, find out about available commands with:
```shell
./kenny -h
```
## TODO
- [x] Integrate with PortAudio for audio recording and audio playback (with the limitation of only 1 concurrent audio stream)
- [x] Use OPUS for audio encoding/decoding
- [ ] Add webRTC signaling client
- [ ] Transmit audio with webRTC
- [ ] Integrate with ion-sfu
- [ ] **Look back and see WTF have I done !?**

## Contribution
Any help and contribution would be greatly appreciated.