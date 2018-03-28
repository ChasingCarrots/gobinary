package gobinary

// StreamReaderView combines a StreamReader with a SeekerView.
// Most of the methods defined here have a corresponding version on the view,
// see there for documentation.
type StreamReaderView struct {
	*StreamReader
	SeekerView
}

func MakeStreamReaderView(reader *StreamReader) StreamReaderView {
	return StreamReaderView{
		StreamReader: reader,
		SeekerView:   MakeSeekerView(reader, 0),
	}
}

func (sr *StreamReaderView) Offset() int64 {
	return sr.SeekerView.Local(sr.GlobalOffset())
}

func (sr *StreamReaderView) GlobalOffset() int64 {
	return sr.StreamReader.Offset()
}

func (sr *StreamReaderView) View(offset int64) {
	sr.SeekerView = sr.SeekerView.View(offset)
}

func (sr *StreamReaderView) ViewHere() {
	sr.SeekerView = sr.SeekerView.View(sr.Offset())
}

func (sv *StreamReaderView) Seek(offset int64, whence int) (int64, error) {
	return sv.SeekerView.Seek(offset, whence)
}
