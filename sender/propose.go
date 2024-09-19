package sender

import (
	"fmt"
	"github.com/6b70/peerbeam/proto/compiled/controlpb"
	"github.com/pion/webrtc/v4"
	"google.golang.org/protobuf/proto"
	"path"
	"time"
)

func (s *Sender) proposeTransfer(ftList []fileTransfer) error {
	fileMDList := &controlpb.FileMetadataList{
		Files: make([]*controlpb.FileMetadata, 0, len(ftList)),
	}
	for _, ft := range ftList {
		fileMDList.Files = append(fileMDList.Files, &controlpb.FileMetadata{
			TransferId:  ft.transferUUID.String(),
			FileName:    path.Base(ft.filePath),
			FileSize:    ft.fileInfo.Size(),
			IsDirectory: ft.fileInfo.IsDir(),
		})
	}

	pbBytes, err := proto.Marshal(fileMDList)
	if err != nil {
		return err
	}
	err = s.DataCh.Send(pbBytes)
	if err != nil {
		return err
	}

	fmt.Println("Waiting for receiver to accept the transfer...")
	err = s.consentCheck()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) consentCheck() error {
	var msg *webrtc.DataChannelMessage
	select {
	case <-time.After(300 * time.Second):
		return fmt.Errorf("timed out waiting for consent")
	case <-s.Ctx.Done():
		return fmt.Errorf("context cancelled")
	case msg = <-s.MsgCh:
		break
	}
	consent := &controlpb.TransferConsent{}
	err := proto.Unmarshal(msg.Data, consent)
	if err != nil {
		return err
	}
	if !consent.Consent {
		return fmt.Errorf("%s", consent.Reason)
	}
	return nil
}
