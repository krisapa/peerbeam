package sender

import (
	"fmt"
	"github.com/6b70/peerbeam/proto/compiled/controlpb"
	"github.com/6b70/peerbeam/utils"
	"google.golang.org/protobuf/proto"
	"path/filepath"
	"time"
)

func (s *Sender) ProposeTransfer(ftList []utils.FileTransfer, answerStr string) error {
	remoteSDP, err := utils.DecodeSDP(answerStr)
	if err != nil {
		return err
	}

	err = s.Session.Conn.SetRemoteDescription(*remoteSDP)
	if err != nil {
		return err
	}

	select {
	case <-s.Session.DataChOpen:
		break
	case <-s.Session.Ctx.Done():
		return fmt.Errorf("context cancelled")
	}

	err = s.sendTransferInfo(ftList)
	if err != nil {
		return err
	}

	err = s.consentCheck()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sender) sendTransferInfo(ftList []utils.FileTransfer) error {
	fileMDList := &controlpb.FileMetadataList{
		Files: make([]*controlpb.FileMetadata, 0, len(ftList)),
	}
	for _, ft := range ftList {
		fileMDList.Files = append(fileMDList.Files, &controlpb.FileMetadata{
			TransferId:  ft.TransferUUID.String(),
			FileName:    filepath.Base(ft.FilePath),
			FileSize:    ft.FileInfo.Size(),
			IsDirectory: ft.FileInfo.IsDir(),
		})
	}

	pbBytes, err := proto.Marshal(fileMDList)
	if err != nil {
		return err
	}
	err = s.Session.DataCh.Send(pbBytes)
	if err != nil {
		return err
	}
	fmt.Println("Sent transfer info")

	return nil
}

func (s *Sender) consentCheck() error {
	msg, err := s.Session.ReceiveMessage(300 * time.Second)
	if err != nil {
		return err
	}

	consent := &controlpb.TransferConsent{}
	err = proto.Unmarshal(msg, consent)
	if err != nil {
		return err
	}
	if !consent.Consent {
		return fmt.Errorf("%s", consent.Reason)
	}
	return nil
}
