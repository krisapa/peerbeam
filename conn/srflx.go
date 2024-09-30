package conn

import (
	"fmt"
	"github.com/pion/webrtc/v4"
	"time"
)

const gatherTimeout = 30 * time.Second

func FetchSRFLX() ([]*webrtc.ICECandidate, error) {
	conn, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{iceServers[0]},
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	res := make([]*webrtc.ICECandidate, 0)
	conn.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i != nil && i.Typ == webrtc.ICECandidateTypeSrflx {
			res = append(res, i)
		}
	})

	done := webrtc.GatheringCompletePromise(conn)
	localSDP, err := conn.CreateOffer(nil)
	if err != nil {
		return nil, err
	}

	err = conn.SetLocalDescription(localSDP)
	if err != nil {
		return nil, err
	}
	select {
	case <-done:
		break
	case <-time.After(gatherTimeout):
		return nil, fmt.Errorf("timed out waiting for ICE gathering to complete")
	}

	return res, nil
}
