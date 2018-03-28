package gobinary

import (
	"io"
)

type StreamReader struct {
	io.ReadSeeker
	offset int64
}

func NewStreamReader(reader io.ReadSeeker) *StreamReader {
	offset, _ := reader.Seek(0, io.SeekCurrent)
	return &StreamReader{
		ReadSeeker: reader,
		offset:     offset,
	}
}

func (sr *StreamReader) Offset() int64 {
	return sr.offset
}

func (sr *StreamReader) Seek(offset int64, whence int) (int64, error) {
	offset, err := sr.ReadSeeker.Seek(offset, whence)
	sr.offset = offset
	return offset, err
}

func (sr *StreamReader) Read(p []byte) (n int, err error) {
	n, err = sr.ReadSeeker.Read(p)
	sr.offset += int64(n)
	return n, err
}

func (sr *StreamReader) SeekCurrent() error {
	_, err := sr.ReadSeeker.Seek(sr.offset, io.SeekStart)
	return err
}
