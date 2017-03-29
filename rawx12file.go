package gox12

import (
	"bufio"
	//"bytes"
	//"fmt"
	"io"
	"log"
	"strings"
)

// X12 line reader and part delimiters
type rawX12FileReader struct {
	reader         *bufio.Reader
	segmentTerm    byte
	elementTerm    byte
	subelementTerm byte
	repetitionTerm byte
	icvn           string
}

type RawSegment struct {
	Segment   Segment
	LineCount int
}

// NewRawX12FileReader creates a RawX12FileReader from a io.Reader.
func NewRawX12FileReader(inFile io.Reader) (*rawX12FileReader, error) {
	const isaLength = 106
	r := new(rawX12FileReader)
	r.reader = bufio.NewReader(inFile)
	//buffer := bytes.NewBuffer(make([]byte, 0))
	first, err := r.reader.Peek(isaLength)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	isa := string(first)
	segTerm, eleTerm, subeleTerm, repTerm, icvn := getDelimiters(isa)
	r.segmentTerm = segTerm
	r.elementTerm = eleTerm
	r.subelementTerm = subeleTerm
	r.repetitionTerm = repTerm
	r.icvn = icvn
	return r, nil
}

func NewRawX12FileReaderFromString(s string) (*rawX12FileReader, error) {
	reader := strings.NewReader(s)
	return NewRawX12FileReader(reader)
}

func (r *rawX12FileReader) GetSegments() <-chan RawSegment {
	ch := make(chan RawSegment)
	ct := 0
	go func() {
		for {
			row, err := r.reader.ReadString(r.segmentTerm)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			if strings.HasSuffix(row, string(r.segmentTerm)) {
				row = row[:len(row)-1]
			}
			row = strings.Trim(row, "\r\n")
			mySeg := NewSegment(row, r.elementTerm, r.subelementTerm, r.repetitionTerm)
			ct++
			seg := RawSegment{
				mySeg,
				ct,
			}
			ch <- seg
		}
		close(ch)
	}()
	return ch
}

// Get the X12 delimiters specified in the ISA segment
func getDelimiters(isa string) (segTerm byte, eleTerm byte, subeleTerm byte, repTerm byte, icvn string) {
	segTerm = isa[len(isa)-1]
	eleTerm = isa[3]
	subeleTerm = isa[len(isa)-2]
	icvn = isa[84:89]
	if icvn == "00501" {
		repTerm = isa[82]
	}
	return
}
