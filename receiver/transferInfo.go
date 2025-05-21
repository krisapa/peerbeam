package receiver

import (
	"fmt"
	"github.com/krisapa/peerbeam/proto/compiled/controlpb"
	"google.golang.org/protobuf/proto"
)

func (r *Receiver) ReceiveTransferInfo() (*controlpb.FileMetadataList, error) {
	select {
	case <-r.Session.DataChOpen:
		break
	case <-r.Session.Ctx.Done():
		return nil, fmt.Errorf("context cancelled")
	}

	dcMSG, err := r.Session.ReceiveMessage(DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var fileMDList controlpb.FileMetadataList
	if err := proto.Unmarshal(dcMSG, &fileMDList); err != nil {
		return nil, err
	}

	return &fileMDList, nil
}

func (r *Receiver) SendTransferConsent(isTransferAccepted bool) error {
	consentBytes, err := proto.Marshal(&controlpb.TransferConsent{
		Consent: isTransferAccepted,
	})
	if err != nil {
		return err
	}
	if err = r.Session.DataCh.Send(consentBytes); err != nil {
		return err
	}
	if !isTransferAccepted {
		// Flushes data channel before exiting
		r.Session.DataCh.GracefulClose()
	}

	return nil
}
