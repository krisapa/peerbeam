package sender

import (
	"fmt"
	"github.com/6b70/peerbeam/conn"
	"github.com/6b70/peerbeam/utils"
	"github.com/pion/webrtc/v4"
)

type Sender struct {
	*conn.Session
}

func New() *Sender {
	return &Sender{
		Session: conn.New(),
	}
}

func (s *Sender) SetupSenderConn() (string, error) {
	err := s.SetupPeerConn()
	if err != nil {
		return "", err
	}
	err = s.setupDataCh()
	if err != nil {
		return "", err
	}

	offer, err := s.createOffer()
	if err != nil {
		return "", err
	}

	return offer, nil
}

func (s *Sender) ProposeTransfer(ftList []utils.FileTransfer, answerStr string) error {
	remoteSDP, err := utils.DecodeSDP(answerStr)
	if err != nil {
		return err
	}

	err = s.AddRemote(remoteSDP)
	if err != nil {
		return err
	}

	select {
	case <-s.DataChOpen:
		break
	case <-s.Ctx.Done():
		return fmt.Errorf("context cancelled")
	}

	err = s.sendTransferInfo(ftList)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) Send(ftList []utils.FileTransfer) error {
	err := s.sendFiles(ftList)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) createOffer() (string, error) {
	initialSDPOffer, err := s.Conn.CreateOffer(nil)
	if err != nil {
		return "", err
	}
	done := webrtc.GatheringCompletePromise(s.Conn)
	err = s.Conn.SetLocalDescription(initialSDPOffer)
	if err != nil {
		return "", err
	}
	<-done
	sdpOffer := s.Conn.LocalDescription()
	encodedSDP, err := utils.EncodeSDP(sdpOffer)
	if err != nil {
		return "", err
	}

	return encodedSDP, nil
}
