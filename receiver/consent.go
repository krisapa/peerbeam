package receiver

import (
	"fmt"
	"github.com/6b70/peerbeam/proto/compiled/controlpb"
	"github.com/6b70/peerbeam/utils"
	"google.golang.org/protobuf/proto"
	"strings"
)

func (r *Receiver) consentToReceive() (*controlpb.FileMetadataList, error) {
	dcMSG := <-r.MsgCh
	if dcMSG.IsString {
		return nil, fmt.Errorf("expected pb, got string")
	}
	var fileMDList controlpb.FileMetadataList
	if err := proto.Unmarshal(dcMSG.Data, &fileMDList); err != nil {
		return nil, err
	}
	for _, fileMD := range fileMDList.Files {
		fmt.Printf("File Name: %s\n", fileMD.FileName)
		fmt.Printf("File Size: %s\n", utils.ByteCountSI(fileMD.FileSize))
		fmt.Printf("Is Directory: %t\n", fileMD.IsDirectory)
	}
	if !inputConsent() {
		return nil, fmt.Errorf("transfer rejected")
	}
	consentBytes, err := proto.Marshal(&controlpb.TransferConsent{
		Consent: true,
	})
	if err != nil {
		return nil, err
	}
	if err = r.DataCh.Send(consentBytes); err != nil {
		return nil, err
	}

	return &fileMDList, nil

}

func inputConsent() bool {
	for {
		fmt.Println("Do you accept the transfer? (y/n)")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(strings.ToLower(input))
		switch input {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("invalid input")
		}
	}
}
