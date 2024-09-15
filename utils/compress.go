package utils

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
)

const pipeBufferSize = 64 * 1024

func CompressStream(file *os.File) *io.PipeReader {
	reader := bufio.NewReader(file)
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		gz := gzip.NewWriter(pw)
		defer gz.Close()

		buf := make([]byte, pipeBufferSize)
		for {
			n, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				pw.CloseWithError(err)
				return
			}
			if n == 0 {
				break
			}
			if _, err = gz.Write(buf[:n]); err != nil {
				pw.CloseWithError(err)
				return
			}
			if err = gz.Flush(); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()
	return pr
}

func DecompressStream(dataReader io.Reader) *io.PipeReader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		gz, err := gzip.NewReader(dataReader)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		defer gz.Close()

		buf := make([]byte, pipeBufferSize)
		for {
			n, err := gz.Read(buf)
			if err != nil && err != io.EOF {
				pw.CloseWithError(err)
				return
			}
			if n == 0 {
				break
			}
			if _, err = pw.Write(buf[:n]); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()
	return pr
}
