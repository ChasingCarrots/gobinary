package gobinary

import (
	"io"
)

// StreamReaderView combines a StreamReader with a SeekerView.
// Most of the methods defined here have a corresponding version on the view,
// see there for documentation.
type StreamReaderView struct {
	*BufferedStreamReader
	view SeekerView
}

func NewStreamReaderView(reader io.ReadSeeker, bufferSize int) StreamReaderView {
	sr := StreamReaderView{}
	sr.BufferedStreamReader = NewBufferedStreamReader(reader, bufferSize)
	sr.view.Init(sr.BufferedStreamReader, 0)
	return sr
}

func (sr *StreamReaderView) Seek(offset int64, whence int) (int64, error) {
	return sr.view.Seek(offset, whence)
}

func (sr *StreamReaderView) Offset() int64 {
	return sr.view.Local(sr.BufferedStreamReader.Offset())
}

func (sr *StreamReaderView) Base() int64 {
	return sr.view.Base()
}

func (sr *StreamReaderView) GlobalOffset() int64 {
	return sr.BufferedStreamReader.Offset()
}

func (sr *StreamReaderView) Local(offset int64) int64 {
	return sr.view.Local(offset)
}

// Copy creates a copy of the the view and its underlying StreamReader.
// This crucially implies that position of this view and the created copy
// are completely independent.
func (sr *StreamReaderView) Copy() StreamReaderView {
	copy := StreamReaderView{}
	copy.BufferedStreamReader = sr.BufferedStreamReader.Copy()
	copy.view.Init(copy.BufferedStreamReader, sr.view.baseOffset)
	return copy
}

func (sr *StreamReaderView) View(offset int64) {
	sr.view = sr.view.View(offset)
}

func (sr *StreamReaderView) ViewHere() {
	sr.view = sr.view.View(sr.Offset())
}
