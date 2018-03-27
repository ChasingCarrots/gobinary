package gobinary

import (
	"errors"
	"io"
)

var NegativeOffset = errors.New("Cannot go to negative offset")

type WriteBuffer struct {
	buffer []byte
	offset int64
}

func (wb *WriteBuffer) Reset() {
	wb.buffer = wb.buffer[:0]
	wb.offset = 0
}

func (wb *WriteBuffer) Offset() int64 { return wb.offset }

func (wb *WriteBuffer) Seek(offset int64, whence int) (int64, error) {
	var target int64
	if whence == io.SeekStart {
		target = offset
	} else if whence == io.SeekEnd {
		target = int64(len(wb.buffer)) + offset
	} else if whence == io.SeekCurrent {
		target = wb.offset + offset
	}
	wb.prepareBuffer(target)
	if target < 0 {
		return wb.offset, NegativeOffset
	}
	wb.offset = target
	return target, nil
}

func (wb *WriteBuffer) prepareBuffer(target int64) {
	if target > int64(cap(wb.buffer)) {
		newBuffer := make([]byte, target, 2*target)
		copy(newBuffer, wb.buffer)
		wb.buffer = newBuffer
	} else if target > int64(len(wb.buffer)) {
		wb.buffer = wb.buffer[:target]
	}
}

func (wb *WriteBuffer) Write(p []byte) (n int, err error) {
	n = len(p)
	target := wb.offset + int64(n)
	wb.prepareBuffer(target)
	buf := wb.buffer[wb.offset:target]
	copy(buf, p)
	wb.offset = target
	return n, nil
}

func (wb *WriteBuffer) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write(wb.buffer)
	return int64(m), err
}

func (wb *WriteBuffer) View(offset int64) WriteBufferView {
	return NewWriteBufferView(wb, offset)
}

func (wb *WriteBuffer) Bytes() []byte {
	n := len(wb.buffer)
	output := make([]byte, n, n)
	copy(output, wb.buffer)
	return output
}
