package gobinary

import "io"

// SeekerView provides a view on a seeker. In this package's terminology, a view
// is simply an offset into the seeker. This allows other abstractions built on
// top of the SeekerView to pretend that they are reading or writing from the
// origin of a stream.
type SeekerView struct {
	io.Seeker
	baseOffset int64
}

// NewSeekerView constructs a new view into a seeker with the given offset.
func NewSeekerView(seeker io.Seeker, offset int64) *SeekerView {
	return &SeekerView{
		Seeker:     seeker,
		baseOffset: offset,
	}
}

// Init initializes a seeker view for the cases where you want to use it as a
// value instead of a pointer type.
func (sv *SeekerView) Init(seeker io.Seeker, offset int64) {
	sv.Seeker = seeker
	sv.baseOffset = offset
}

// Seek seeks in the underlying seeker. When whence is SeekCurrent or SeekEnd,
// this is the same as seeking in the underlying seeker. For whence equal to
// SeekStart however, this is seeking relative to the base of this view.
func (sv *SeekerView) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		offset += sv.baseOffset
	}
	o, err := sv.Seeker.Seek(offset, whence)
	return o - sv.baseOffset, err
}

// Offset returns the local offset that the underlying seeker is currently at.
// For example, if the underlying seeker is at the base of this view, Offset is
// 0. Note that the value returned by this function may be negative when the
// underlying seeker is currently at a position smaller than this view's base.
func (sv *SeekerView) Offset() int64 {
	return sv.Local(sv.GlobalOffset())
}

// Base returns the real offset of this view in the underlying seeker. This is
// the global offset that this view considers its local 0 offset.
func (sv *SeekerView) Base() int64 {
	return sv.baseOffset
}

// GlobalOffset returns the global (real, absolute) offset of the underlying seeker.
func (sv *SeekerView) GlobalOffset() int64 {
	offset, _ := sv.Seeker.Seek(0, io.SeekCurrent)
	return offset
}

// Local takes an offset in global (real, absolute) units and returns its offset from
// this view.
func (sv *SeekerView) Local(absOffset int64) int64 {
	return absOffset - sv.baseOffset
}

// View creates a new view at the given local offset. Note that the underlying seeker
// of the new view is the same as the underlying seeker of this view.
func (sv *SeekerView) View(offset int64) SeekerView {
	return SeekerView{
		Seeker:     sv.Seeker,
		baseOffset: sv.baseOffset + offset,
	}
}

// ViewHere creates a new view at the current offset of the underlying seeker.
// Note that the underlying seeker of the new view is the same as the underlying seeker
// of this view.
func (sv *SeekerView) ViewHere() SeekerView {
	return SeekerView{
		Seeker:     sv.Seeker,
		baseOffset: sv.baseOffset + sv.Offset(),
	}
}
