package sender

import (
	"fmt"
	"github.com/6b70/peerbeam/conn"
	"github.com/6b70/peerbeam/utils"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v4"
	"os"
	"path/filepath"
)

type Sender struct {
	*conn.Session
}

func StartSender(files []string) error {
	sender := New()
	return sender.Send(files)
}

func New() *Sender {
	return &Sender{
		Session: conn.New(),
	}
}

type fileTransfer struct {
	transferUUID uuid.UUID
	filePath     string
	fileInfo     os.FileInfo
}

func (s *Sender) Send(files []string) error {
	defer s.CtxCancel()
	ftList, err := parseFiles(files)
	if err != nil {
		return err
	}

	err = s.SetupPeerConn()
	if err != nil {
		return err
	}

	err = s.setupDataCh()
	if err != nil {
		return err
	}

	offer, err := s.CreateOffer()
	if err != nil {
		return err
	}

	utils.CopyGeneratedSDPPrompt(offer)
	fmt.Println("1. Offer copied to clipboard. Share it with the receiver.")
	fmt.Println("2. Copy the receiver's answer and press Enter.")
	remoteSDP := utils.InputSDPPrompt()
	err = s.AddRemote(remoteSDP)
	if err != nil {
		return err
	}

	<-s.DataChOpen
	err = s.proposeTransfer(ftList)
	if err != nil {
		return err
	}

	err = s.sendFiles(ftList)
	if err != nil {
		return err
	}
	fmt.Println("Files sent successfully")
	return nil
}

func (s *Sender) CreateOffer() (string, error) {
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

func parseFiles(files []string) ([]fileTransfer, error) {
	ftList := make([]fileTransfer, 0, len(files))
	for _, relFP := range files {
		fp, err := filepath.Abs(relFP)
		if err != nil {
			return nil, fmt.Errorf("error with file '%s': %v", relFP, err)
		}
		fi, err := os.Stat(fp)
		if err != nil {
			return nil, fmt.Errorf("error with file '%s': %v", fp, err)
		}
		ftList = append(ftList, fileTransfer{
			transferUUID: uuid.New(),
			filePath:     fp,
			fileInfo:     fi,
		})
	}
	return ftList, nil
}
