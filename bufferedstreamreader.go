package gobinary

import (
	"fmt"
	"io"
)

// BufferedStreamReader provides helper methods for reading common basic types from an underlying
// stream of data. The stream is defined by the characteristic that by default, each
// reading operation proceeds forward.
// You can also seek within the stream.
// In essence, the BufferedStreamReader should be thought of as a point within the stream (its current
// offset) that is automatically advanced whenever a reading operation is requested.
type BufferedStreamReader struct {
	HighLevelReader
	reader io.ReadSeeker
	// The position within the current buffer
	bufPos int
	// The buffer used to buffer the reading.
	buf []byte
	// A temporary buffer for reading large structs and reading across buffer limits.
	tmp []byte
	// The actual position in the underlying stream that the buffer was read from.
	offset     int64
	bufferSize int
}

// NewBufferedStreamReader produces a new stream reader for the underlying reading. The buffer
// size specifies the size of the buffer used to prevent frequent reading from the underlying
// stream.
func NewBufferedStreamReader(reader io.ReadSeeker, bufferSize int) *BufferedStreamReader {
	sr := BufferedStreamReader{}
	sr.init(reader, bufferSize)
	sr.offset, _ = reader.Seek(0, io.SeekCurrent)
	return &sr
}

func (sr *BufferedStreamReader) init(reader io.ReadSeeker, bufferSize int) {
	const minSize = 16
	if bufferSize < minSize {
		bufferSize = minSize
	}
	sr.reader = reader
	sr.buf = make([]byte, 0, bufferSize)
	sr.tmp = make([]byte, 0, minSize)
	sr.bufPos = 0
	sr.bufferSize = bufferSize
	sr.HighLevelReader.ByteReader = sr
}

// Copy returns a new BufferedStreamReader that for all intents and purposes is a perfect copy
// of this BufferedStreamReader, but has a offset that is decoupled from this reader's offset.
// You can have multiple readers into the same underlying reader; but using them is not
// thread-safe.
func (sr *BufferedStreamReader) Copy() *BufferedStreamReader {
	cp := BufferedStreamReader{}
	cp.init(sr.reader, cap(sr.buf))
	cp.bufPos = sr.bufPos
	cp.offset = sr.offset
	cp.buf = make([]byte, len(sr.buf), cap(sr.buf))
	copy(cp.buf, sr.buf)
	return &cp
}

// Offset returns the offset of the reader in the underlying reader. This is the position
// at which the next read will be issued.
func (sr *BufferedStreamReader) Offset() int64 {
	return sr.offset + int64(sr.bufPos)
}

// Seek moves the offset of this reader to the given offset. This does not seek on the underlying
// stream until a read is requested.
func (sr *BufferedStreamReader) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		sr.seekImpl(offset)
		return sr.Offset(), nil
	} else if whence == io.SeekCurrent {
		sr.seekImpl(sr.Offset() + offset)
		return sr.Offset(), nil
	} else {
		offset, err := sr.reader.Seek(offset, whence)
		sr.seekImpl(offset)
		return sr.Offset(), err
	}
}

func (sr *BufferedStreamReader) seekImpl(offset int64) {
	delta := offset - (sr.offset + int64(sr.bufPos))
	sr.bufPos += int(delta)
	if sr.bufPos < 0 || sr.bufPos > len(sr.buf) {
		sr.InvalidateBuffer()
		sr.offset = offset
	}
}

// Read is here to implement the io.Reader interface and work as described there.
func (sr *BufferedStreamReader) Read(p []byte) (n int, err error) {
	n = sr.ReadBytes(p)
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

// BufferSize returns the original buffer size specified when this BufferedStreamReader was
// created.
func (sr *BufferedStreamReader) BufferSize() int { return sr.bufferSize }

// InvalidateBuffer invalidates the underlying buffer, which forces any reading
// operation to read from the underlying stream. Use this only when you absolutely
// need to ensure that the data read is the most recently written data.
func (sr *BufferedStreamReader) InvalidateBuffer() {
	sr.buf = sr.buf[0:0]
	sr.bufPos = 0
}

// bigRead is to be used for big reads that will definitely go
// over the size of one or multiple buffers.
func (sr *BufferedStreamReader) bigRead(bytes int) []byte {
	if cap(sr.tmp) < bytes {
		sr.tmp = make([]byte, bytes, 2*bytes+1)
	} else {
		sr.tmp = sr.tmp[:bytes]
	}

	output := sr.tmp
	copy(output, sr.buf[sr.bufPos:len(sr.buf)])
	read := len(sr.buf) - sr.bufPos
	bytes -= read
	for bytes > 0 {
		end := cap(sr.buf)
		if bytes < end {
			end = bytes
		}

		sr.offset = sr.offset + int64(len(sr.buf))
		sr.reader.Seek(sr.offset, io.SeekStart)
		sr.buf = sr.buf[:cap(sr.buf)]
		n, err := sr.reader.Read(sr.buf)

		if n == 0 && err != nil || n < end {
			panic(fmt.Sprintf("unexpected end of stream: %v, read %v, wanted %v", err, n, len(sr.buf)))
		}
		sr.buf = sr.buf[0:n]
		bytes -= copy(output, sr.buf[0:end])
		sr.bufPos = end
	}
	return sr.tmp
}

// smallRead is to be used for reads that will definitely fit into the
// buffer. Note that due to their position they could still require us to read
// into the next buffer.
func (sr *BufferedStreamReader) smallRead(bytes int) []byte {
	oldPos := sr.bufPos
	sr.bufPos += bytes
	if sr.bufPos <= len(sr.buf) {
		return sr.buf[oldPos:sr.bufPos]
	}
	sr.tmp = sr.tmp[:bytes]
	copied := copy(sr.tmp, sr.buf[oldPos:len(sr.buf)])

	sr.offset += int64(len(sr.buf))
	sr.reader.Seek(sr.offset, io.SeekStart)
	sr.buf = sr.buf[:cap(sr.buf)]
	n, err := sr.reader.Read(sr.buf)

	remaining := bytes - copied
	if n == 0 && err != nil || n < remaining {
		panic(fmt.Sprintf("unexpected end of stream: %v, read %v, wanted %v", err, n, len(sr.buf)))
	}
	sr.buf = sr.buf[:n]
	copy(sr.tmp[copied:bytes], sr.buf)
	return sr.tmp
}

// GetReadBuffer implements the ByteReader interface.
func (sr *BufferedStreamReader) GetReadBuffer(bytes int) []byte {
	if bytes < cap(sr.buf) {
		return sr.smallRead(bytes)
	}
	return sr.bigRead(bytes)
}
