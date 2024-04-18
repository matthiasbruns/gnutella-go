package whatsmyip

import (
	"log"

	"github.com/pion/webrtc/v3"
)

// DiscoverPublicIP discovers public IP address of executed device by STUN server
func DiscoverPublicIP(cb func(string, error)) {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		log.Println("err connection")
		cb("", err)
		return
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			cb("", nil)
		}

		// recieve public ip address
		if c != nil && c.Typ == webrtc.ICECandidateTypeSrflx {
			cb(c.Address, nil)
		}
	})

	if _, err := peerConnection.CreateDataChannel("", nil); err != nil {
		log.Println("err crerate data channel")
		cb("", err)
		return
	}

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Println("err crerate offer")
		cb("", err)
		return
	}

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		log.Println("err set local description")
		cb("", err)
		return
	}
}
