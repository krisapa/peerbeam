package receiver

import (
	"fmt"
	"github.com/ksp237/peerbeam/conn"
	"github.com/ksp237/peerbeam/utils"
	"github.com/pion/webrtc/v4"
)

type Receiver struct {
	*conn.Session
}

func New() *Receiver {
	return &Receiver{
		Session: conn.New(),
	}
}

func (r *Receiver) ReceiveFiles(destPath string) error {
	defer r.CtxCancel()
	err := r.SetupPeerConn()
	if err != nil {
		return err
	}

	candidatePromise := webrtc.GatheringCompletePromise(r.Conn)
	remoteSDP := utils.InputSDPPrompt()
	r.addChHandler()
	err = r.AddRemote(remoteSDP)
	if err != nil {
		return err
	}
	answer, err := r.CreateAnswer(candidatePromise)
	if err != nil {
		return err
	}
	fmt.Println("Copy the answer and send it to the sender:")
	fmt.Println(answer)
	<-r.DataChOpen
	fileMDList, err := r.consentToReceive()
	if err != nil {
		return err
	}
	err = r.receiveFiles(fileMDList, destPath)
	if err != nil {
		return err
	}

	return nil
}

func (r *Receiver) CreateAnswer(candidatePromise <-chan struct{}) (string, error) {
	initAnswer, err := r.Conn.CreateAnswer(nil)
	if err != nil {
		return "", err
	}
	err = r.Conn.SetLocalDescription(initAnswer)
	if err != nil {
		return "", err
	}
	<-candidatePromise
	answer := r.Conn.LocalDescription()

	encodedSDP, err := utils.EncodeSDP(answer)
	if err != nil {
		return "", err
	}

	return encodedSDP, nil
}
