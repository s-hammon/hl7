package hl7

import "sync"

type scanner struct {
	err error

	fldDelim byte
	comDelim byte
	repDelim byte
	escDelim byte
	subDelim byte
}

var scannerPool = sync.Pool{
	New: func() any {
		return &scanner{}
	},
}

func newScanner() *scanner {
	s := scannerPool.Get().(*scanner)
	s.err = nil
	return s
}

func (s *scanner) state(c byte) int {
	switch c {
	default:
		return stateValue
	case s.fldDelim:
		return stateFieldIdx
	// TODO: add cases for other delimiters
	case '\r':
		return stateEndSegment
	}
}
