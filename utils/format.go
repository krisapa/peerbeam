package utils

import (
	"fmt"
	"github.com/krisapa/peerbeam/proto/compiled/controlpb"
	"strings"
)

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func FormatFileProposal(fileMDList *controlpb.FileMetadataList) string {
	var builder strings.Builder
	for _, fileMD := range fileMDList.Files {
		name := fileMD.FileName
		if fileMD.IsDirectory {
			name += "/"
		}
		size := ByteCountSI(fileMD.FileSize)
		builder.WriteString(fmt.Sprintf("%s\t%s\n", name, size))
	}
	return builder.String()
}
