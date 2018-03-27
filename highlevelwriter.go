package gobinary

import (
	"encoding/binary"
	"math"
)

type ByteWriter interface {
	// GetWriteBuffer should return a write buffer that has enough space
	// to hold the next bytes. This buffer is assumed to be temporary and
	// should not be kept by any caller.
	GetWriteBuffer(length int) []byte
}

type HighLevelWriter struct {
	ByteWriter
}

func (hlw HighLevelWriter) WriteInt64(value int64) {
	binary.LittleEndian.PutUint64(hlw.GetWriteBuffer(8), uint64(value))
}

func (hlw HighLevelWriter) WriteInt32(value int32) {
	binary.LittleEndian.PutUint32(hlw.GetWriteBuffer(4), uint32(value))
}

func (hlw HighLevelWriter) WriteInt16(value int16) {
	binary.LittleEndian.PutUint16(hlw.GetWriteBuffer(2), uint16(value))
}

func (hlw HighLevelWriter) WriteInt8(value int8) {
	hlw.GetWriteBuffer(1)[0] = byte(value)
}

func (hlw HighLevelWriter) WriteUInt64(value uint64) {
	binary.LittleEndian.PutUint64(hlw.GetWriteBuffer(8), value)
}

func (hlw HighLevelWriter) WriteUInt32(value uint32) {
	binary.LittleEndian.PutUint32(hlw.GetWriteBuffer(4), value)
}

func (hlw HighLevelWriter) WriteUInt16(value uint16) {
	binary.LittleEndian.PutUint16(hlw.GetWriteBuffer(2), value)
}

func (hlw HighLevelWriter) WriteUInt8(value uint8) {
	hlw.GetWriteBuffer(1)[0] = byte(value)
}

func (hlw HighLevelWriter) WriteByte(value byte) {
	hlw.GetWriteBuffer(1)[0] = value
}

func (hlw HighLevelWriter) WriteFloat32(value float32) {
	hlw.WriteUInt32(math.Float32bits(value))
}

func (hlw HighLevelWriter) WriteFloat64(value float64) {
	hlw.WriteUInt64(math.Float64bits(value))
}

func (hlw HighLevelWriter) WriteBool(value bool) {
	b := uint8(0)
	if value {
		b = 1
	}
	hlw.WriteUInt8(b)
}

func (hlw HighLevelWriter) WriteString(value string) {
	copy(hlw.GetWriteBuffer(len(value)), value)
}

func (hlw HighLevelWriter) WriteBytes(value []byte) {
	copy(hlw.GetWriteBuffer(len(value)), value)
}
