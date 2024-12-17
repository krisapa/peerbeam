package utils

import (
	"github.com/klauspost/compress/zstd"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func CompressStream(file *os.File, blockSize int) *io.PipeReader {
	pr, pw := io.Pipe()
	enc, err := zstd.NewWriter(pw)
	if err != nil {
		pw.CloseWithError(err)
		return nil
	}
	buf := make([]byte, blockSize)

	go func() {
		defer func() {
			err := enc.Close()
			if err != nil {
				pw.CloseWithError(err)
			} else {
				pw.Close()
			}
		}()
		_, err = io.CopyBuffer(enc, file, buf)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
	}()
	return pr
}

func DecompressStream(dataReader io.Reader, blockSize int) *io.PipeReader {
	pr, pw := io.Pipe()
	d, err := zstd.NewReader(dataReader)
	if err != nil {
		pw.CloseWithError(err)
		return nil
	}
	buf := make([]byte, blockSize)

	go func() {
		defer func() {
			d.Close()
			pw.Close()
		}()
		_, err = io.CopyBuffer(pw, d, buf)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
	}()
	return pr
}

var archiveExtensions = []string{".zip", ".tar", ".gz", ".rar", ".7z"}

func IsArchiveFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return slices.Contains(archiveExtensions, ext)
}
