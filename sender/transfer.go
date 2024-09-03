package sender

import (
	"fmt"
	"github.com/ksp237/peerbeam/proto/compiled/transferpb"
	"github.com/pion/webrtc/v4"
	"google.golang.org/protobuf/proto"
	"os"
	"path"
)

func (s *Sender) sendFiles(ftList []fileTransfer) error {
	for _, ft := range ftList {
		err := s.sendFile(ft)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sender) sendFile(ft fileTransfer) error {
	fmt.Printf("Sending file '%s'...\n", path.Base(ft.filePath))

	fileBytes, err := os.ReadFile(ft.filePath)
	if err != nil {
		return err
	}

	transferPB := &transferpb.File{
		TransferId: ft.transferUUID.String(),
		Data:       fileBytes,
		IsCompress: true,
		IsEncrypt:  false,
	}

	pbBytes, err := proto.Marshal(transferPB)
	if err != nil {
		return err
	}

	err = s.DataCh.Send(pbBytes)
	if err != nil {
		return err
	}
	var fileResp *webrtc.DataChannelMessage
	select {
	case <-s.Ctx.Done():
		return fmt.Errorf("context cancelled")
	case fileResp = <-s.MsgCh:
		break
	}
	var fileRespPB transferpb.TransferComplete
	if err = proto.Unmarshal(fileResp.Data, &fileRespPB); err != nil {
		return err
	}
	if fileRespPB.TransferId != ft.transferUUID.String() {
		return fmt.Errorf("unexpected transfer ID: %s", fileRespPB.TransferId)
	}
	if !fileRespPB.Success {
		return fmt.Errorf("transfer failed: %s", fileRespPB.Message)
	}

	fmt.Printf("File '%s' sent successfully\n", path.Base(ft.filePath))
	return nil
}
