package gobinary

import (
	"encoding/binary"
	"io"
	"math"
)

type HighLevelReader struct {
	io.Reader
	buffer []byte
}

func MakeHighLevelReader(reader io.Reader) HighLevelReader {
	return HighLevelReader{
		Reader: reader,
		buffer: make([]byte, 0, 16),
	}
}

func (hlr *HighLevelReader) getReadBuffer(length int) []byte {
	if cap(hlr.buffer) < length {
		c := 2*length + 1
		if c < 16 {
			c = 16
		}
		hlr.buffer = make([]byte, length, c)
	} else {
		hlr.buffer = hlr.buffer[0:length]
	}
	hlr.Reader.Read(hlr.buffer)
	return hlr.buffer
}

// ReadInt64 reads a 64bit integer from the stream.
func (hlr *HighLevelReader) ReadInt64() int64 {
	return int64(binary.LittleEndian.Uint64(hlr.getReadBuffer(8)))
}

// ReadInt32 reads a 32bit integer from the stream.
func (hlr *HighLevelReader) ReadInt32() int32 {
	return int32(binary.LittleEndian.Uint32(hlr.getReadBuffer(4)))
}

// ReadInt16 reads a 16bit integer from the stream.
func (hlr *HighLevelReader) ReadInt16() int16 {
	return int16(binary.LittleEndian.Uint16(hlr.getReadBuffer(2)))
}

// ReadInt8 reads an 8bit integer from the stream.
func (hlr *HighLevelReader) ReadInt8() int8 {
	return int8(hlr.getReadBuffer(1)[0])
}

// ReadUInt64 reads a 64bit unsigned integer from the stream.
func (hlr *HighLevelReader) ReadUInt64() uint64 {
	return binary.LittleEndian.Uint64(hlr.getReadBuffer(8))
}

// ReadUInt32 reads a 32bit unsigned integer from the stream.
func (hlr *HighLevelReader) ReadUInt32() uint32 {
	return binary.LittleEndian.Uint32(hlr.getReadBuffer(4))
}

// ReadUInt16 reads a 16bit unsigned integer from the stream.
func (hlr *HighLevelReader) ReadUInt16() uint16 {
	return binary.LittleEndian.Uint16(hlr.getReadBuffer(2))
}

// ReadUInt8 reads a 8bit unsigned integer from the stream.
func (hlr *HighLevelReader) ReadUInt8() uint8 {
	return uint8(hlr.getReadBuffer(1)[0])
}

// ReadByte reads a single byte from the underlying stream.
func (hlr *HighLevelReader) ReadByte() byte {
	return hlr.getReadBuffer(1)[0]
}

// ReadFloat32 reads a 32bit floating point value in IEEE754 format
// from the underlying stream.
func (hlr *HighLevelReader) ReadFloat32() float32 {
	return math.Float32frombits(hlr.ReadUInt32())
}

// ReadFloat64 reads a 64bit floating point value in IEEE754 format
// from the underlying stream.
func (hlr *HighLevelReader) ReadFloat64() float64 {
	return math.Float64frombits(hlr.ReadUInt64())
}

// ReadBool reads a boolean value from the underlying stream. Booleans
// are encoded as single byte with a value of 0 encoding false.
func (hlr *HighLevelReader) ReadBool() bool {
	return hlr.getReadBuffer(1)[0] > 0
}

// ReadString reads a UTF8 encoded string of the given number of bytes (not runes!)
// from the underlying stream.
func (hlr *HighLevelReader) ReadString(length int) string {
	return string(hlr.getReadBuffer(length))
}

// ReadBytes reads as many bytes as the passed in byte slice is long into said byte
// slice. Returns the number of bytes read.
func (hlr *HighLevelReader) ReadBytes(output []byte) int {
	l := len(output)
	tmp := hlr.getReadBuffer(l)
	copy(output, tmp)
	return len(tmp)
}
