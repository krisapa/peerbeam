package sender

import "github.com/6b70/peerbeam/utils"

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
	sdpOffer := s.Session.Conn.LocalDescription()
	s.Session.WaitGatherComplete()
	encodedSDP, err := utils.EncodeSDP(sdpOffer, s.Session.Candidates)
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
	return nil
}
