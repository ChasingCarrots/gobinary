package gobinary

type StreamWriterView struct {
	*StreamWriter
	SeekerView
}

func MakeStreamWriterView(writer *StreamWriter) StreamWriterView {
	sw := StreamWriterView{}
	sw.StreamWriter = writer
	sw.SeekerView.Init(sw.StreamWriter, 0)
	return sw
}

func (sw *StreamWriterView) Offset() int64 {
	return sw.SeekerView.Local(sw.GlobalOffset())
}

func (sw *StreamWriterView) GlobalOffset() int64 {
	return sw.StreamWriter.Offset()
}

func (sw *StreamWriterView) View(offset int64) {
	sw.SeekerView = sw.SeekerView.View(offset)
}

func (sw *StreamWriterView) ViewHere() {
	sw.SeekerView = sw.SeekerView.View(sw.Offset())
}

func (sv *StreamWriterView) Seek(offset int64, whence int) (int64, error) {
	return sv.SeekerView.Seek(offset, whence)
}
