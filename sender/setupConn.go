package sender

import (
	"github.com/6b70/peerbeam/utils"
	"github.com/pion/webrtc/v4"
)

func (s *Sender) SetupSenderConn() (string, error) {
	err := s.Session.SetupPeerConn()
	if err != nil {
		return "", err
	}
	err = s.createChs()
	if err != nil {
		return "", err
	}

	offer, err := s.createOffer()
	if err != nil {
		return "", err
	}

	return offer, nil
}

func (s *Sender) createOffer() (string, error) {
	initialSDPOffer, err := s.Session.Conn.CreateOffer(nil)
	if err != nil {
		return "", err
	}
	err = s.Session.Conn.SetLocalDescription(initialSDPOffer)
	if err != nil {
		return "", err
	}

	gatherComplete := make(chan struct{})
	s.Session.Conn.OnICEGatheringStateChange(func(state webrtc.ICEGatheringState) {
		if state == webrtc.ICEGatheringStateComplete {
			close(gatherComplete)
		}
	})
	<-gatherComplete

	sdpOffer := s.Session.Conn.LocalDescription()
	encodedSDP, err := utils.EncodeSDP(sdpOffer)
	if err != nil {
		return "", err
	}

	return encodedSDP, nil
}

func (s *Sender) createChs() error {
	ch, err := s.Session.Conn.CreateDataChannel("data", nil)
	if err != nil {
		return err
	}
	s.Session.DataChHandler(ch)

	ch, err = s.Session.Conn.CreateDataChannel("candidate", nil)
	if err != nil {
		return err
	}
	s.Session.CandidateChHandler(ch)

	return nil
}
