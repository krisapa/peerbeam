package receiver

import (
	"fmt"
	"github.com/6b70/peerbeam/proto/compiled/controlpb"
	"google.golang.org/protobuf/proto"
)

func (r *Receiver) ReceiveFileProposal() (*controlpb.FileMetadataList, error) {
	select {
	case <-r.DataChOpen:
		break
	case <-r.Ctx.Done():
		return nil, fmt.Errorf("context cancelled")
	}

	dcMSG := <-r.MsgCh
	var fileMDList controlpb.FileMetadataList
	if err := proto.Unmarshal(dcMSG.Data, &fileMDList); err != nil {
		return nil, err
	}
	return &fileMDList, nil
}

func (r *Receiver) SendProposalResponse(isTransferAccepted bool) error {
	consentBytes, err := proto.Marshal(&controlpb.TransferConsent{
		Consent: isTransferAccepted,
	})
	if err != nil {
		return err
	}
	if err = r.DataCh.Send(consentBytes); err != nil {
		return err
	}
	return nil
}
