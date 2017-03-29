package gox12

import "bytes"

type X12File struct {
	Segments []Segment
}

func (f *X12File) AppendSegment(s Segment) {
	f.Segments = append(f.Segments, s)
}

func (f X12File) String() string {
	var buf bytes.Buffer
	for _, s := range f.Segments {
		buf.WriteString(s.String())
		buf.WriteString("\n")
	}
	return buf.String()
}
