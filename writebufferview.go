package gobinary

import "io"

type WriteBufferView struct {
	buffer     *WriteBuffer
	baseOffset int64
}

func NewWriteBufferView(buffer *WriteBuffer, baseOffset int64) WriteBufferView {
	return WriteBufferView{
		buffer:     buffer,
		baseOffset: baseOffset,
	}
}

func (view *WriteBufferView) Offset() int64 {
	return view.buffer.Offset() - view.baseOffset
}

func (view *WriteBufferView) Write(p []byte) (n int, err error) {
	return view.buffer.Write(p)
}

func (view *WriteBufferView) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		offset += view.baseOffset
	}
	off, err := view.buffer.Seek(offset, whence)
	off -= view.baseOffset
	return off, err
}

func (view *WriteBufferView) View(offset int64) WriteBufferView {
	return NewWriteBufferView(view.buffer, view.baseOffset+offset)
}
