package hl7

import (
	"fmt"
	"testing"

	"github.com/s-hammon/p"
	"github.com/stretchr/testify/require"
)

type sillyParser struct {
	d      *decodeState
	hl7Idx int

	result map[string]map[int]string
}

func newParser(data []byte) *sillyParser {
	d := &decodeState{}
	d.init(data)
	return &sillyParser{d: d, result: make(map[string]map[int]string)}
}

func (p *sillyParser) Scan() bool {
	if p.d.savedError != nil || p.d.prev == stateEOF {
		return false
	}
	p.d.scanValue()

	return true
}

func (p *sillyParser) ReadVal(start int) string {
	return string(p.d.data[start : p.d.off-1])
}

func (p *sillyParser) newSegment(segment string) {
	if _, ok := p.result[segment]; !ok {
		p.result[segment] = make(map[int]string)
	}
}

func (p *sillyParser) addField(segment string, val string) {
	p.result[segment][p.hl7Idx] = val
}

func TestDecodeState_Read(t *testing.T) {
	t.Parallel()

	parser := newParser([]byte("MSH|^~\\&|SendingApp|SendingFac|ReceivingApp|ReceivingFac|20250101120000||ADT^A01|1234567|P|2.3\rEVN|A01|20250101120000|"))

	segment := "MSH"
	parser.newSegment(segment)
	parser.hl7Idx = 1
	parser.addField(segment, string(parser.d.scan.fldDelim))
	parser.hl7Idx = 2
	parser.addField(segment, parser.d.encodingChars())

	start := parser.d.off
	for parser.Scan() {
		val := parser.ReadVal(start)
		if val != "" {
			parser.addField(segment, val)
			t.Logf("segment: %q, idx: %d, value: %q\n", segment, parser.hl7Idx, val)
		}

		if parser.d.prev == stateEndSegment {
			t.Log("new segment state")

			start := parser.d.off
			parser.d.scanN(3)
			if parser.d.off > len(parser.d.data) {
				break
			}

			segment = string(parser.d.data[start:parser.d.off])
			if len(segment) != 3 {
				parser.d.savedError = fmt.Errorf("malformed segment name %q", segment)
				break
			}

			parser.newSegment(segment)
			parser.hl7Idx = 0
		} else {
			parser.hl7Idx++
		}

		start = parser.d.off
	}

	require.NoError(t, parser.d.savedError)
	require.Contains(t, p.Keys(parser.result), "EVN")
	t.Log(parser.result["MSH"])
	require.Equal(t, "A01", parser.result["EVN"][1])
}
