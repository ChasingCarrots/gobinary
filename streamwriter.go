package gobinary

import (
	"io"
)

// StreamWriter conceptually represents an offset in the Stream.
// Calling one of the writing method advances the offset.
type StreamWriter struct {
	io.WriteSeeker
	// The offset in the underlying stream.
	offset int64
}

// NewStreamWriter creates a new StreamWriter atop the given writer.
func NewStreamWriter(writer io.WriteSeeker) *StreamWriter {
	sw := StreamWriter{
		WriteSeeker: writer,
	}
	sw.offset, _ = sw.WriteSeeker.Seek(0, io.SeekCurrent)
	return &sw
}

// Offset returns the offset in the underlying stream at which the next
// writing operation will be issued.
func (sw *StreamWriter) Offset() int64 {
	return sw.offset
}

// Seek seeks in the underlying writer, flushing the buffer of this writer
// if necessary.
func (sw *StreamWriter) Seek(offset int64, whence int) (n int64, err error) {
	sw.offset, err = sw.WriteSeeker.Seek(offset, whence)
	return sw.offset, err
}

func (sw *StreamWriter) Write(p []byte) (n int, err error) {
	n, err = sw.WriteSeeker.Write(p)
	sw.offset += int64(n)
	return n, err
}

func (sw *StreamWriter) SeekCurrent() error {
	_, err := sw.WriteSeeker.Seek(sw.offset, io.SeekStart)
	return err
}
