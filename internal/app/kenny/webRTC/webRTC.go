package webRTC

import (
	"github.com/pion/webrtc/v3"
)

// make config for stun & turn servers (no turn server yet :))
func makeConfig() webrtc.Configuration {
	config := webrtc.Configuration {
		ICEServers: []webrtc.ICEServer {
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	return config
}

// return a new peerConnection
func makePeerConnection() *webrtc.PeerConnection{
	peerConnection, err := webrtc.NewPeerConnection(makeConfig())
	if err != nil {
		panic(err)
	}
	return peerConnection
}