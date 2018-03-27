package gobinary

import "io"

type StreamWriterView struct {
	*BufferedStreamWriter
	view SeekerView
}

func NewStreamWriterView(writer io.WriteSeeker, bufferSize int) StreamWriterView {
	sw := StreamWriterView{}
	sw.BufferedStreamWriter = NewBufferedStreamWriter(writer, bufferSize)
	sw.view.Init(sw.BufferedStreamWriter, 0)
	return sw
}

func (sw *StreamWriterView) Seek(offset int64, whence int) (int64, error) {
	return sw.view.Seek(offset, whence)
}

func (sw *StreamWriterView) Offset() int64 {
	return sw.view.Local(sw.BufferedStreamWriter.Offset())
}

func (sw *StreamWriterView) Base() int64 {
	return sw.view.Base()
}

func (sw *StreamWriterView) GlobalOffset() int64 {
	return sw.BufferedStreamWriter.Offset()
}

func (sw *StreamWriterView) Local(absOffset int64) int64 {
	return sw.view.Local(absOffset)
}

func (sw *StreamWriterView) Copy() StreamWriterView {
	sw.Flush()
	copy := StreamWriterView{}
	copy.BufferedStreamWriter = NewBufferedStreamWriter(sw.BufferedStreamWriter.writer, sw.BufferedStreamWriter.BufferSize())
	copy.view.Init(copy.BufferedStreamWriter, 0)
	return copy
}

func (sw *StreamWriterView) View(offset int64) {
	sw.view = sw.view.View(offset)
}

func (sw *StreamWriterView) ViewHere() {
	sw.view = sw.view.View(sw.Offset())
}
