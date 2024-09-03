package receiver

import (
	"fmt"
	"github.com/ksp237/peerbeam/proto/compiled/controlpb"
	"github.com/ksp237/peerbeam/proto/compiled/transferpb"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"time"
)

func (r *Receiver) receiveFiles(fileMDList *controlpb.FileMetadataList, destPath string) error {
	transferIDMap := make(map[string]*controlpb.FileMetadata)
	for _, fileMD := range fileMDList.Files {
		transferIDMap[fileMD.TransferId] = fileMD
	}
	for i := 0; i < len(fileMDList.Files); i++ {
		var filePB transferpb.File
		select {
		case dcMSG := <-r.MsgCh:
			if dcMSG.IsString {
				return fmt.Errorf("expected pb, got string")
			}
			if err := proto.Unmarshal(dcMSG.Data, &filePB); err != nil {
				return err
			}
		case <-r.Ctx.Done():
			return fmt.Errorf("context cancelled")
		}
		fmd, ok := transferIDMap[filePB.TransferId]
		if !ok {
			return fmt.Errorf("unexpected file ID received: %s", filePB.TransferId)
		}
		if err := handleFile(&filePB, fmd, destPath); err != nil {
			return err
		}

		fmt.Printf("Received file '%s'\n", fmd.FileName)
		fmt.Printf("File size: %d\n", len(filePB.Data))

		pbBytes, err := proto.Marshal(&transferpb.TransferComplete{
			TransferId: filePB.TransferId,
			Success:    true,
		})
		if err != nil {
			return err
		}
		if err = r.DataCh.Send(pbBytes); err != nil {
			return err
		}
	}
	// flush last message
	time.Sleep(1 * time.Second)

	return nil
}

func handleFile(filePB *transferpb.File, fmd *controlpb.FileMetadata, destPath string) error {
	fileBytes := filePB.Data
	filePath := filepath.Join(destPath, fmd.FileName)
	return os.WriteFile(filePath, fileBytes, 0644)
}
