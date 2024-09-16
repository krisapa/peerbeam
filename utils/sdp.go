package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v4"
	"golang.design/x/clipboard"
	"io"
)

func EncodeSDP(sdp *webrtc.SessionDescription) (string, error) {
	sdpJSON, err := json.Marshal(sdp)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	g, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	defer g.Close()
	if _, err = g.Write(sdpJSON); err != nil {
		return "", err
	}

	if err = g.Close(); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeSDP(in string) (*webrtc.SessionDescription, error) {
	buf, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	sdpBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var sdp webrtc.SessionDescription
	err = json.Unmarshal(sdpBytes, &sdp)
	if err != nil {
		return nil, err
	}

	return &sdp, nil
}

func InputSDPPrompt() *webrtc.SessionDescription {
	for {
		fmt.Scanln()
		clipBytes := clipboard.Read(clipboard.FmtText)
		if remoteSDP, err := DecodeSDP(string(clipBytes)); err == nil {
			fmt.Println("SDP successfully read.")
			return remoteSDP
		}
		fmt.Println("Invalid SDP. Ensure it's copied to your clipboard and press Enter.")
	}
}

func CopyGeneratedSDPPrompt(sdp string) {
	clipboard.Write(clipboard.FmtText, []byte(sdp))
}
