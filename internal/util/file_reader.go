package util

import (
	"io"
	"os"
)

type FilerReader struct {
	*os.File
}

func (fr FilerReader) Read(p []byte) (n int, err error) {
	_, err = fr.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return fr.File.Read(p)
}
