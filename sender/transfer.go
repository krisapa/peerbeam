package sender

import (
	"bufio"
	"fmt"
	"github.com/6b70/peerbeam/proto/compiled/transferpb"
	"github.com/6b70/peerbeam/utils"
	"github.com/pion/webrtc/v4"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"time"
)

const (
	blockSize                  int    = 32 * 1024
	bufferedAmountLowThreshold uint64 = 512 * 1024  // 512 KB
	maxBufferedAmount          uint64 = 1024 * 1024 // 1 MB
)

func (s *Sender) sendFiles(ftList []utils.FileTransfer) error {
	s.DataCh.SetBufferedAmountLowThreshold(bufferedAmountLowThreshold)

	sendMoreCh := make(chan struct{}, 1)
	s.DataCh.OnBufferedAmountLow(func() {
		select {
		case sendMoreCh <- struct{}{}:
		default:
		}
	})

	for i, ft := range ftList {
		err := s.sendFile(ft, sendMoreCh)
		if err != nil {
			return err
		}
		if err = s.transferConfirmation(ft); err != nil {
			if i == len(ftList)-1 {
				return nil
			}
			return err
		}
	}
	return nil
}

func (s *Sender) sendFile(ft utils.FileTransfer, sendMoreCh <-chan struct{}) error {
	isCompressed := !utils.IsArchiveFile(ft.FilePath)

	if err := s.sendTransferStart(ft, isCompressed); err != nil {
		return err
	}
	if err := s.sendFileBlocks(ft, sendMoreCh, isCompressed); err != nil {
		return err
	}
	return nil
}

func (s *Sender) sendTransferStart(ft utils.FileTransfer, isCompressed bool) error {
	transferStartBytes, err := proto.Marshal(&transferpb.TransferStart{
		TransferId:   ft.TransferUUID.String(),
		IsCompressed: isCompressed,
	})
	if err != nil {
		return err
	}
	return s.DataCh.Send(transferStartBytes)
}

func (s *Sender) sendFileBlocks(ft utils.FileTransfer, sendMoreCh <-chan struct{}, isCompressed bool) error {
	peerInfoStr, err := s.PeerInfoStr()
	if err != nil {
		return err
	}
	bar := utils.NewProgressBar(ft.FileInfo.Size(), peerInfoStr, true)
	defer bar.Close()

	file, err := os.Open(ft.FilePath)
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
		err = bar.Add(n)
		if err != nil {
			bar.Finish()
		}
	}
	return nil
}

func (s *Sender) sendFileBlock(ft utils.FileTransfer, data []byte, isLastBlock bool, sendMoreCh <-chan struct{}) error {
	fileBlock := &transferpb.FileBlock{
		TransferId:  ft.TransferUUID.String(),
		Data:        data,
		IsLastBlock: isLastBlock,
	}

	pbBytes, err := proto.Marshal(fileBlock)
	if err != nil {
		return err
	}

	if s.DataCh.BufferedAmount() >= maxBufferedAmount {
		select {
		case <-sendMoreCh:
		case <-s.Ctx.Done():
			return fmt.Errorf("context cancelled while waiting to send data")
		}
	}

	return s.DataCh.Send(pbBytes)
}

func (s *Sender) transferConfirmation(ft utils.FileTransfer) error {
	var fileResp *webrtc.DataChannelMessage
	select {
	case <-s.Ctx.Done():
		return fmt.Errorf("context cancelled while waiting for confirmation response")
	case <-time.After(1 * time.Second):
		return fmt.Errorf("timeout while waiting for confirmation response")
	case fileResp = <-s.MsgCh:
	}

	var fileRespPB transferpb.TransferComplete
	if err := proto.Unmarshal(fileResp.Data, &fileRespPB); err != nil {
		return err
	}
	if fileRespPB.TransferId != ft.TransferUUID.String() {
		return fmt.Errorf("unexpected transfer ID: %s", fileRespPB.TransferId)
	}
	if !fileRespPB.Success {
		return fmt.Errorf("transfer failed: %s", fileRespPB.Message)
	}
	return nil
}
