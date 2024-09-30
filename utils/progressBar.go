package utils

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
)

func NewProgressBar(totalBytes int64, peerInfoStr string, isSender bool) *progressbar.ProgressBar {
	var descStr string
	if isSender {
		descStr = fmt.Sprintf("Sending (->%s)", peerInfoStr)
	} else {
		descStr = fmt.Sprintf("Receiving (<-%s)", peerInfoStr)
	}
	return progressbar.DefaultBytes(totalBytes, descStr)
}
