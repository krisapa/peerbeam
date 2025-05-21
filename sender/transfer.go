package sender

import (
	"bufio"
	"fmt"
	"github.com/krisapa/peerbeam/proto/compiled/transferpb"
	"github.com/krisapa/peerbeam/utils"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"time"
)

const (
	// BlockSize max size for a single packet is 65535 for pion/sctp
	BlockSize                         = 65000
	maxBufferedAmount          uint64 = 1024 * 1024 // 1 MB
	bufferedAmountLowThreshold        = maxBufferedAmount - BlockSize
)

func (s *Sender) SendFiles(ftList []utils.FileTransfer) error {
	s.Session.DataCh.SetBufferedAmountLowThreshold(bufferedAmountLowThreshold)

	sendMoreCh := make(chan struct{}, 1)
	s.Session.DataCh.OnBufferedAmountLow(func() {
		select {
		case sendMoreCh <- struct{}{}:
		default:
		}
	})

	for i, ft := range ftList {
		isCompressed := !utils.IsArchiveFile(ft.FilePath)
		if err := s.sendTransferStart(ft, isCompressed); err != nil {
			return err
		}
		if err := s.sendFileBlocks(ft, isCompressed, sendMoreCh); err != nil {
			return err
		}
		if err := s.transferConfirmation(ft); err != nil {
			if i == len(ftList)-1 {
				return nil
			}
			return err
		}
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
	return s.Session.DataCh.Send(transferStartBytes)
}

func (s *Sender) sendFileBlocks(ft utils.FileTransfer, isCompressed bool, sendMoreCh <-chan struct{}) error {
	peerInfoStr, err := s.Session.PeerInfoStr()
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
		reader = utils.CompressStream(file, BlockSize)
	} else {
		reader = bufio.NewReader(file)
	}

	isLastBlock := false
	for !isLastBlock {
		fileBytes := make([]byte, BlockSize)
		n, err := reader.Read(fileBytes)
		if err != nil && err != io.EOF {
			return err
		}
		isLastBlock = err == io.EOF
		if err = s.sendFileBlock(ft, fileBytes[:n], isLastBlock, sendMoreCh); err != nil {
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

	if s.Session.DataCh.BufferedAmount() >= maxBufferedAmount {
		select {
		case <-sendMoreCh:
		case <-s.Session.Ctx.Done():
			return fmt.Errorf("context cancelled while waiting to send data")
		}
	}

	return s.Session.DataCh.Send(pbBytes)
}

func (s *Sender) transferConfirmation(ft utils.FileTransfer) error {
	fileRespBytes, err := s.Session.ReceiveMessage(7 * time.Second)
	if err != nil {
		return err
	}

	var fileRespPB transferpb.TransferComplete
	if err = proto.Unmarshal(fileRespBytes, &fileRespPB); err != nil {
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
