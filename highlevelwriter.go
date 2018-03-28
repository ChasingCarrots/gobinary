package gobinary

import (
	"encoding/binary"
	"io"
	"math"
)

type HighLevelWriter struct {
	io.Writer
	buffer []byte
}

func MakeHighLevelWriter(writer io.Writer) HighLevelWriter {
	return HighLevelWriter{
		Writer: writer,
		buffer: make([]byte, 0, 16),
	}
}

func (hlw *HighLevelWriter) getWriteBuffer(length int) []byte {
	if cap(hlw.buffer) < length {
		c := 2*length + 1
		if c < 16 {
			c = 16
		}
		hlw.buffer = make([]byte, length, c)
	} else {
		hlw.buffer = hlw.buffer[0:length]
	}
	return hlw.buffer
}

func (hlw *HighLevelWriter) commitWrite() (int, error) {
	return hlw.Writer.Write(hlw.buffer)
}

func (hlw *HighLevelWriter) WriteInt64(value int64) (int, error) {
	binary.LittleEndian.PutUint64(hlw.getWriteBuffer(8), uint64(value))
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteInt32(value int32) (int, error) {
	binary.LittleEndian.PutUint32(hlw.getWriteBuffer(4), uint32(value))
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteInt16(value int16) (int, error) {
	binary.LittleEndian.PutUint16(hlw.getWriteBuffer(2), uint16(value))
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteInt8(value int8) (int, error) {
	hlw.getWriteBuffer(1)[0] = byte(value)
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteUInt64(value uint64) (int, error) {
	binary.LittleEndian.PutUint64(hlw.getWriteBuffer(8), value)
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteUInt32(value uint32) (int, error) {
	binary.LittleEndian.PutUint32(hlw.getWriteBuffer(4), value)
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteUInt16(value uint16) (int, error) {
	binary.LittleEndian.PutUint16(hlw.getWriteBuffer(2), value)
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteUInt8(value uint8) (int, error) {
	hlw.getWriteBuffer(1)[0] = byte(value)
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteByte(value byte) (int, error) {
	hlw.getWriteBuffer(1)[0] = value
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteFloat32(value float32) (int, error) {
	return hlw.WriteUInt32(math.Float32bits(value))
}

func (hlw *HighLevelWriter) WriteFloat64(value float64) (int, error) {
	return hlw.WriteUInt64(math.Float64bits(value))
}

func (hlw *HighLevelWriter) WriteBool(value bool) (int, error) {
	b := uint8(0)
	if value {
		b = 1
	}
	return hlw.WriteUInt8(b)
}

func (hlw *HighLevelWriter) WriteString(value string) (int, error) {
	copy(hlw.getWriteBuffer(len(value)), value)
	return hlw.commitWrite()
}

func (hlw *HighLevelWriter) WriteBytes(value []byte) (int, error) {
	_, err := hlw.Writer.Write(value)
	return 0, err
}
