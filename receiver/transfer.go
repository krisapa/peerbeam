package receiver

import (
	"fmt"
	"github.com/6b70/peerbeam/proto/compiled/controlpb"
	"github.com/6b70/peerbeam/proto/compiled/transferpb"
	"github.com/6b70/peerbeam/utils"
	"github.com/pion/webrtc/v4"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path/filepath"
)

func (r *Receiver) receiveFiles(fileMDList *controlpb.FileMetadataList, destPath string) error {
	transferIDMap := make(map[string]*controlpb.FileMetadata)
	for _, fileMD := range fileMDList.Files {
		transferIDMap[fileMD.TransferId] = fileMD
	}

	for _ = range fileMDList.Files {
		transferStartMsg, err := r.receiveBlock()
		if err != nil {
			return err
		}
		var transferStart transferpb.TransferStart
		if err = proto.Unmarshal(transferStartMsg.Data, &transferStart); err != nil {
			return err
		}
		fmd, ok := transferIDMap[transferStart.TransferId]
		if !ok {
			return fmt.Errorf("unexpected transfer ID received: %s", transferStart.TransferId)
		}
		err = r.receiveFile(&transferStart, fmd, destPath)
		if err != nil {
			return err
		}
		fmt.Printf("Received file '%s'\n", fmd.FileName)

		// Send the confirmation response
		pbBytes, err := proto.Marshal(&transferpb.TransferComplete{
			TransferId: transferStart.TransferId,
			Success:    true,
		})
		if err != nil {
			return err
		}
		if err = r.DataCh.Send(pbBytes); err != nil {
			return err
		}
	}

	return nil
}

func (r *Receiver) receiveBlock() (*webrtc.DataChannelMessage, error) {
	select {
	case dcMSG := <-r.MsgCh:
		return dcMSG, nil
	case <-r.Ctx.Done():
		return nil, fmt.Errorf("context cancelled while receiving file")
	}
}

func (r *Receiver) receiveFile(ts *transferpb.TransferStart, fmd *controlpb.FileMetadata, destPath string) error {
	filePath := filepath.Join(destPath, fmd.FileName)
	destFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	decompressedReader, err := r.receiveAndDecompress(ts)
	if err != nil {
		return err
	}
	if _, err = io.Copy(destFile, decompressedReader); err != nil {
		return err
	}
	return nil
}

func (r *Receiver) receiveAndDecompress(ts *transferpb.TransferStart) (io.Reader, error) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		isLastBlock := false
		for !isLastBlock {
			msg, err := r.receiveBlock()
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			var fileBlock transferpb.FileBlock
			if err = proto.Unmarshal(msg.Data, &fileBlock); err != nil {
				pw.CloseWithError(err)
				return
			}
			if fileBlock.TransferId != ts.TransferId {
				pw.CloseWithError(fmt.Errorf("unexpected transfer ID received: %s", fileBlock.TransferId))
				return
			}
			isLastBlock = fileBlock.IsLastBlock
			if _, err = pw.Write(fileBlock.Data); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()
	return utils.DecompressStream(pr), nil
}
