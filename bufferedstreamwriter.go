package gobinary

import (
	"fmt"
	"io"
)

// BufferedStreamWriter provides helper methods for writing common basic types to an
// underlying writer.
// A BufferedStreamWriter conceptually represents an offset in the Stream plus a buffer
// for buffered writing. Calling one of the writing method advances the offset
// and may or may not flush the buffer to the underlying writer.
type BufferedStreamWriter struct {
	HighLevelWriter
	buffer []byte
	writer io.WriteSeeker
	// The offset in the underlying stream where the current buffer should be
	// written to.
	offset     int64
	bufferSize int
}

// NewBufferedStreamWriter creates a new BufferedStreamWriter atop the given writer with a
// given size for the buffer used in buffered writing. Please note that the
// buffer size cannot be zero, it has a minimum value of 16.
func NewBufferedStreamWriter(writer io.WriteSeeker, bufferSize int) *BufferedStreamWriter {
	sw := BufferedStreamWriter{}
	sw.init(writer, bufferSize)
	return &sw
}

// BufferSize returns the original buffer size used when this writer was
// created.
func (sw *BufferedStreamWriter) BufferSize() int {
	return sw.bufferSize
}

func (sw *BufferedStreamWriter) init(writer io.WriteSeeker, bufferSize int) {
	if bufferSize < 16 {
		bufferSize = 16
	}
	sw.bufferSize = bufferSize
	sw.buffer = make([]byte, 0, bufferSize)
	sw.writer = writer
	sw.offset, _ = sw.writer.Seek(0, io.SeekCurrent)
	sw.HighLevelWriter.ByteWriter = sw
}

// Offset returns the offset in the underlying stream at which the next
// writing operation will be issued.
func (sw *BufferedStreamWriter) Offset() int64 {
	return sw.offset + int64(len(sw.buffer))
}

// Seek seeks in the underlying writer, flushing the buffer of this writer
// if necessary.
func (sw *BufferedStreamWriter) Seek(offset int64, whence int) (n int64, err error) {
	if (whence == io.SeekStart && offset == sw.Offset()) ||
		(whence == io.SeekCurrent && offset == 0) {
		return offset, nil
	}
	sw.Flush()
	sw.offset, err = sw.writer.Seek(offset, whence)
	return sw.offset, err
}

// Flush forces the buffer of this stream to be flushed to the underlying writer.
func (sw *BufferedStreamWriter) Flush() {
	if len(sw.buffer) == 0 {
		return
	}
	sw.writer.Seek(sw.offset, io.SeekStart)
	n, err := sw.writer.Write(sw.buffer)
	if err != nil || n < len(sw.buffer) {
		panic(fmt.Sprintf("Writing to underlying stream failed: %v", err))
	}
	sw.offset += int64(n)
	sw.buffer = sw.buffer[:0]
}

// GetWriteBuffer returns a temporary buffer to write n bytes to.
func (sw *BufferedStreamWriter) GetWriteBuffer(n int) []byte {
	c := cap(sw.buffer)
	l := len(sw.buffer)
	if n <= c-l {
		sw.buffer = sw.buffer[:l+n]
		return sw.buffer[l : l+n]
	} else if n < c {
		sw.Flush()
		sw.buffer = sw.buffer[0:n]
		return sw.buffer
	} else {
		sw.Flush()
		sw.buffer = make([]byte, n, 2*n+1)
		return sw.buffer
	}
}
