package webRTC

import (
	"github.com/pion/webrtc/v3"
)

// Make config for stun & turn servers (no turn server yet :))
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

// Return a new peer connection
func makePeerConnection(onICEStateChaneCallback func(connectionState webrtc.ICEConnectionState)) *webrtc.PeerConnection {
	peerConnection, err := webrtc.NewPeerConnection(makeConfig())
	if err != nil {
		panic(err)
	}
	
	peerConnection.OnICEConnectionStateChange(onICEStateChaneCallback)

	return peerConnection
}

// Take a offer from another system and set it for this peer connection
func setSDPOffer(pc *webrtc.PeerConnection, offer webrtc.SessionDescription) {
	err := pc.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}
}

// After an offer received, we set an answer for that peer connection
func setSDPAnswer(pc *webrtc.PeerConnection) {
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	err = pc.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}
}