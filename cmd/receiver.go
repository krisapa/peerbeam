package cmd

import (
	"fmt"
	"github.com/6b70/peerbeam/proto/compiled/controlpb"
	"github.com/6b70/peerbeam/receiver"
	"github.com/6b70/peerbeam/utils"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"os"
	"path/filepath"
)

func startReceiver(destPath string) error {
	var err error
	destPath, err = filepath.Abs(destPath)
	if err != nil {
		return err
	}
	destInfo, err := os.Stat(destPath)
	if err != nil {
		return err
	}
	if !destInfo.IsDir() {
		return fmt.Errorf("destination path must be a directory")
	}
	r := receiver.New()
	defer r.CtxCancel()

	err = r.SetupReceiverConn()
	if err != nil {
		return err
	}
	offer, err := recvExchangeSDP()
	if err != nil {
		return err
	}
	var (
		actionErr error
		answer    string
	)
	err = spinner.New().
		Type(spinner.Dots).
		Title("Generating answer").
		Action(func() {
			answer, actionErr = r.CreateAnswer(offer)
		}).
		Run()
	if err != nil {
		return err
	}
	if actionErr != nil {
		return actionErr
	}

	utils.CopyGeneratedSDPPrompt(answer)

	var fileMDList *controlpb.FileMetadataList
	err = spinner.New().
		Type(spinner.Dots).
		Title("Answer copied. Send to sender.").
		Action(func() {
			fileMDList, actionErr = r.ReceiveFileProposal()
		}).
		Run()
	if err != nil {
		return err
	}
	if actionErr != nil {
		return actionErr
	}

	fileProposalStr := utils.FormatFileProposal(fileMDList)
	isTransferAccepted, err := transferConsent(fileProposalStr)
	if err != nil {
		return err
	}

	err = r.SendProposalResponse(isTransferAccepted)
	if err != nil {
		return err
	}
	if !isTransferAccepted {
		// Flushes data channel before exiting
		r.DataCh.GracefulClose()
		return nil
	}

	return r.Receive(fileMDList, destPath)
}

func recvExchangeSDP() (string, error) {
	var offer string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Enter sender's offer").
				CharLimit(5000).
				Placeholder("Paste offer here...").
				Value(&offer).
				Validate(utils.ValidateSDP),
		),
	)
	err := form.Run()
	if err != nil {
		return "", err
	}

	return offer, nil
}

func transferConsent(proposalStr string) (bool, error) {
	var isTransferAccepted bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Accept files?").
				Description(proposalStr).
				Affirmative("Yes!").
				Negative("No.").
				Value(&isTransferAccepted),
		),
	)

	err := form.Run()
	if err != nil {
		return false, err
	}
	return isTransferAccepted, nil
}
