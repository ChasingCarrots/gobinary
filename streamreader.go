package gobinary

import (
	"fmt"
	"io"
)

type StreamReader struct {
	reader     io.ReadSeeker
	readBuffer []byte
	offset     int64
}

func NewStreamReader(reader io.ReadSeeker) *StreamReader {
	offset, _ := reader.Seek(0, io.SeekCurrent)
	return &StreamReader{
		reader: reader,
		offset: offset,
	}
}

func (sr *StreamReader) Offset() int64 {
	return sr.offset
}

func (sr *StreamReader) Seek(offset int64, whence int) (int64, error) {
	offset, err := sr.reader.Seek(offset, whence)
	sr.offset = offset
	return offset, err
}

func (sr *StreamReader) Read(p []byte) (n int, err error) {
	sr.reader.Seek(sr.offset, io.SeekStart)
	n, err = sr.reader.Read(p)
	sr.offset += int64(n)
	return n, err
}

func (sr *StreamReader) GetReadBuffer(bytes int) []byte {
	if bytes > cap(sr.readBuffer) {
		sr.readBuffer = make([]byte, bytes, 2*bytes+1)
	} else {
		sr.readBuffer = sr.readBuffer[:bytes]
	}
	_, err := sr.Read(sr.readBuffer)
	if err != nil {
		panic(fmt.Sprintf("Read failed: %v", err))
	}
	return sr.readBuffer
}
