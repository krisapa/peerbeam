package sender

import (
	"bufio"
	"fmt"
	"github.com/ksp237/peerbeam/proto/compiled/transferpb"
	"github.com/ksp237/peerbeam/utils"
	"github.com/pion/webrtc/v4"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path"
)

const (
	blockSize                  int    = 32 * 1024
	bufferedAmountLowThreshold uint64 = 512 * 1024  // 512 KB
	maxBufferedAmount          uint64 = 1024 * 1024 // 1 MB
)

func (s *Sender) sendFiles(ftList []fileTransfer) error {
	s.DataCh.SetBufferedAmountLowThreshold(bufferedAmountLowThreshold)

	sendMoreCh := make(chan struct{}, 1)
	s.DataCh.OnBufferedAmountLow(func() {
		select {
		case sendMoreCh <- struct{}{}:
		default:
		}
	})

	for _, ft := range ftList {
		err := s.sendFile(ft, sendMoreCh)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sender) sendFile(ft fileTransfer, sendMoreCh <-chan struct{}) error {
	isCompressed := !utils.IsArchiveFile(ft.filePath)

	fmt.Printf("Sending file '%s'...\n", path.Base(ft.filePath))
	if err := s.sendTransferStart(ft, isCompressed); err != nil {
		return err
	}
	if err := s.sendFileBlocks(ft, sendMoreCh, isCompressed); err != nil {
		return err
	}

	if err := s.transferConfirmation(ft); err != nil {
		return err
	}

	fmt.Printf("File '%s' sent successfully\n", path.Base(ft.filePath))
	return nil
}

// Helper function to send transfer start message
func (s *Sender) sendTransferStart(ft fileTransfer, isCompressed bool) error {
	transferStartBytes, err := proto.Marshal(&transferpb.TransferStart{
		TransferId:   ft.transferUUID.String(),
		IsCompressed: isCompressed,
	})
	if err != nil {
		return err
	}
	return s.DataCh.Send(transferStartBytes)
}

func (s *Sender) sendFileBlocks(ft fileTransfer, sendMoreCh <-chan struct{}, isCompressed bool) error {
	file, err := os.Open(ft.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader
	if isCompressed {
		reader = utils.CompressStream(file)
	} else {
		reader = bufio.NewReader(file)
	}

	isLastBlock := false
	for !isLastBlock {
		fileBytes := make([]byte, blockSize)
		n, err := reader.Read(fileBytes)
		if err != nil && err != io.EOF {
			return err
		}

		isLastBlock = err == io.EOF
		if err := s.sendFileBlock(ft, fileBytes[:n], isLastBlock, sendMoreCh); err != nil {
			return err
		}
	}
	return nil
}

// Helper function to send individual file block
func (s *Sender) sendFileBlock(ft fileTransfer, data []byte, isLastBlock bool, sendMoreCh <-chan struct{}) error {
	fileBlock := &transferpb.FileBlock{
		TransferId:  ft.transferUUID.String(),
		Data:        data,
		IsLastBlock: isLastBlock,
	}

	pbBytes, err := proto.Marshal(fileBlock)
	if err != nil {
		return err
	}

	// Wait for buffered amount to be reduced
	if s.DataCh.BufferedAmount() >= maxBufferedAmount {
		select {
		case <-sendMoreCh:
		case <-s.Ctx.Done():
			return fmt.Errorf("context cancelled while waiting to send data")
		}
	}

	return s.DataCh.Send(pbBytes)
}

func (s *Sender) transferConfirmation(ft fileTransfer) error {
	var fileResp *webrtc.DataChannelMessage
	select {
	case <-s.Ctx.Done():
		return fmt.Errorf("context cancelled while waiting for confirmation response")
	case fileResp = <-s.MsgCh:
	}

	var fileRespPB transferpb.TransferComplete
	if err := proto.Unmarshal(fileResp.Data, &fileRespPB); err != nil {
		return err
	}
	if fileRespPB.TransferId != ft.transferUUID.String() {
		return fmt.Errorf("unexpected transfer ID: %s", fileRespPB.TransferId)
	}
	if !fileRespPB.Success {
		return fmt.Errorf("transfer failed: %s", fileRespPB.Message)
	}
	return nil
}

//func (s *Sender) sendFile(ft fileTransfer, sendMoreCh <-chan struct{}) error {
//	fmt.Printf("Sending file '%s'...\n", path.Base(ft.filePath))
//
//	transferStartBytes, err := proto.Marshal(&transferpb.TransferStart{
//		TransferId:   ft.transferUUID.String(),
//		IsCompressed: true,
//	})
//	if err != nil {
//		return err
//	}
//	err = s.DataCh.Send(transferStartBytes)
//	if err != nil {
//		return err
//	}
//
//	file, err := os.Open(ft.filePath)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	reader := bufio.NewReader(file)
//
//	isLastBlock := false
//	for !isLastBlock {
//		fileBytes := make([]byte, blockSize)
//		n, err := reader.Read(fileBytes)
//		if err != nil && err != io.EOF {
//			return err
//		}
//
//		isLastBlock = err == io.EOF
//		fileBlock := &transferpb.FileBlock{
//			TransferId:  ft.transferUUID.String(),
//			Data:        fileBytes[:n],
//			IsLastBlock: isLastBlock,
//		}
//
//		pbBytes, err := proto.Marshal(fileBlock)
//		if err != nil {
//			return err
//		}
//
//		if s.DataCh.BufferedAmount() >= maxBufferedAmount {
//			select {
//			case <-sendMoreCh:
//			case <-s.Ctx.Done():
//				return fmt.Errorf("context cancelled while waiting to send data")
//			}
//		}
//
//		err = s.DataCh.Send(pbBytes)
//		if err != nil {
//			return err
//		}
//	}
//
//	var fileResp *webrtc.DataChannelMessage
//	select {
//	case <-s.Ctx.Done():
//		return fmt.Errorf("context cancelled while waiting for confirmation response")
//	case fileResp = <-s.MsgCh:
//		break
//	}
//
//	var fileRespPB transferpb.TransferComplete
//	if err = proto.Unmarshal(fileResp.Data, &fileRespPB); err != nil {
//		return err
//	}
//	if fileRespPB.TransferId != ft.transferUUID.String() {
//		return fmt.Errorf("unexpected transfer ID: %s", fileRespPB.TransferId)
//	}
//	if !fileRespPB.Success {
//		return fmt.Errorf("transfer failed: %s", fileRespPB.Message)
//	}
//
//	fmt.Printf("File '%s' sent successfully\n", path.Base(ft.filePath))
//	return nil
//}

//func (s *Sender) sendFile(ft fileTransfer, sendMoreCh <-chan struct{}) error {
//	fmt.Printf("Sending file '%s'...\n", path.Base(ft.filePath))
//
//	fileBytes, err := os.ReadFile(ft.filePath)
//	if err != nil {
//		return err
//	}
//
//	fileBlock := &transferpb.FileBlock{
//		TransferId:  ft.transferUUID.String(),
//		Data:        fileBytes,
//		IsLastBlock: true,
//	}
//
//	pbBytes, err := proto.Marshal(transferPB)
//	if err != nil {
//		return err
//	}
//
//	err = s.DataCh.Send(pbBytes)
//	if err != nil {
//		return err
//	}
//	var fileResp *webrtc.DataChannelMessage
//	select {
//	case <-s.Ctx.Done():
//		return fmt.Errorf("context cancelled")
//	case fileResp = <-s.MsgCh:
//		break
//	}
//	var fileRespPB transferpb.TransferComplete
//	if err = proto.Unmarshal(fileResp.Data, &fileRespPB); err != nil {
//		return err
//	}
//	if fileRespPB.TransferId != ft.transferUUID.String() {
//		return fmt.Errorf("unexpected transfer ID: %s", fileRespPB.TransferId)
//	}
//	if !fileRespPB.Success {
//		return fmt.Errorf("transfer failed: %s", fileRespPB.Message)
//	}
//
//	fmt.Printf("File '%s' sent successfully\n", path.Base(ft.filePath))
//	return nil
//}
