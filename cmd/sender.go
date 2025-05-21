package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/krisapa/peerbeam/sender"
	"github.com/krisapa/peerbeam/utils"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

func startSender(files []string) error {
	s := sender.New()
	defer s.Session.CtxCancel()

	offerCh := make(chan string)
	go func() {
		offer, err := s.SetupSenderConn()
		if err != nil {
			log.Errorln(err)
			close(offerCh)
		}
		offerCh <- offer
	}()

	ftList, err := utils.ParseFiles(files)
	if err != nil {
		return err
	}

	var offer string
	err = spinner.New().
		Type(spinner.Dots).
		Title("Creating offer").
		Action(func() {
			offer = <-offerCh
		}).
		Run()
	if err != nil {
		return err
	}
	if offer == "" {
		return fmt.Errorf("failed to create offer")
	}

	answer, err := senderExchangeSDP(offer)
	if err != nil {
		return err
	}

	var isTransferAccepted atomic.Bool
	err = spinner.New().
		Type(spinner.Dots).
		Title("Waiting for receiver to accept transfer").
		Action(func() {
			err = s.ProposeTransfer(ftList, answer)
			isTransferAccepted.Store(err == nil)
		}).
		Run()
	if err != nil {
		return err
	}

	if !isTransferAccepted.Load() {
		return fmt.Errorf("transfer rejected")
	}

	return s.SendFiles(ftList)
}

// exchangeSDP handles the SDP offer and answer exchange process
func senderExchangeSDP(offer string) (string, error) {
	utils.CopyGeneratedSDPPrompt(offer)

	var answer string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Offer copied. Send to receiver.").
				CharLimit(5000).
				Placeholder("Paste response here...").
				Value(&answer).
				Validate(func(s string) error {
					// here is no need to check
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		return "", err
	}

	return answer, nil
}
