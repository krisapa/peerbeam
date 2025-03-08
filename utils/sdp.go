package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/aymanbagabas/go-osc52/v2"
	"github.com/pion/webrtc/v4"
	"io"
	"os"
)

type SwapInfo struct {
	Sdp        *webrtc.SessionDescription `json:"sdp"`
	Candidates []webrtc.ICECandidateInit  `json:"candidates"`
}

func EncodeSDP(sdp *webrtc.SessionDescription, candidates []webrtc.ICECandidateInit) (string, error) {
	info := SwapInfo{
		Sdp:        sdp,
		Candidates: candidates,
	}

	infoJSON, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	g, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	defer g.Close()
	if _, err = g.Write(infoJSON); err != nil {
		return "", err
	}

	if err = g.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeSDP(in string) (*webrtc.SessionDescription, []webrtc.ICECandidateInit, error) {
	buf, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()

	infoBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	var info SwapInfo
	err = json.Unmarshal(infoBytes, &info)
	if err != nil {
		return nil, nil, err
	}
	return info.Sdp, info.Candidates, nil
}

func CopyGeneratedSDPPrompt(sdp string) {
	clipboard.WriteAll(sdp)
	fmt.Println("You can manually copy the following content:")
	fmt.Println(sdp)
	osc52.New(sdp).WriteTo(os.Stderr)
}
