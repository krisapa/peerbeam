package receiver

import (
	"fmt"
	"github.com/6b70/peerbeam/conn"
	"github.com/6b70/peerbeam/utils"
	"github.com/pion/webrtc/v4"
)

type Receiver struct {
	*conn.Session
}

func StartReceiver(destPath string) error {
	receiver := New()
	return receiver.Receive(destPath)
}

func New() *Receiver {
	return &Receiver{
		Session: conn.New(),
	}
}

func (r *Receiver) Receive(destPath string) error {
	defer r.CtxCancel()
	err := r.SetupPeerConn()
	if err != nil {
		return err
	}

	candidatePromise := webrtc.GatheringCompletePromise(r.Conn)
	fmt.Println("Copy the sender's offer to the clipboard and press enter.")
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

	utils.CopyGeneratedSDPPrompt(answer)
	fmt.Println("Answer copied to clipboard. Send the answer to the sender.")

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
