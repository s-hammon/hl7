package hl7

import (
	"errors"
	"fmt"
	"testing"

	"github.com/s-hammon/p"
	"github.com/stretchr/testify/require"
)

func TestDecodeState_Unmarshal(t *testing.T) {
	t.Parallel()

	d := decodeState{}
	d.init([]byte("MSH|^~\\&|SendingApp|SendingFac|ReceivingApp|ReceivingFac|20250101120000||ADT^A01|1234567|P|2.3\rEVN|A01|20250101120000|"))

	err := d.unmarshal()
	require.NoError(t, err)
	require.Equal(t, byte('|'), d.fldDelim)
	require.Equal(t, byte('^'), d.comDelim)
	require.Equal(t, byte('~'), d.repDelim)
	require.Equal(t, byte('\\'), d.escDelim)
	require.Equal(t, byte('&'), d.subDelim)
	require.Equal(t, 3, d.hl7FieldIdx)
}

type sillyParser struct {
	d    *decodeState
	last int
}

func newParser(data []byte) *sillyParser {
	d := &decodeState{}
	d.init(data)
	return &sillyParser{d: d, last: d.off}
}

func (p *sillyParser) Scan() bool {
	if p.d.savedError != nil || p.d.state == stateEOF {
		return false
	}
	p.last = p.d.off
	if err := p.d.readNext(); err != nil {
		return false
	}

	return true
}

func (p *sillyParser) ReadVal() string {
	i, j := p.last, p.d.off-1
	return string(p.d.data[i:j])
}

func TestDecodeState_Read(t *testing.T) {
	t.Parallel()

	parser := newParser([]byte("MSH|^~\\&|SendingApp|SendingFac|ReceivingApp|ReceivingFac|20250101120000||ADT^A01|1234567|P|2.3\rEVN|A01|20250101120000|"))

	segs := make(map[string]map[int]string)
	segment := "MSH"

	segs[segment] = make(map[int]string)
	segs[segment][1] = string(parser.d.fldDelim)
	segs[segment][2] = parser.d.encodingChars()

	last := parser.d.state
	i := parser.d.hl7FieldIdx
	for parser.Scan() {
		if last == stateSegmentName {
			t.Log("new segment state")
			val := parser.ReadVal()
			if val == "" {
				parser.d.savedError = errors.New("expecting segment name...")
			} else if len(val) != 3 {
				parser.d.savedError = fmt.Errorf("malformed segment name %q", val)
			}

			segment = val
		} else {
			val := parser.ReadVal()
			if val != "" {
				seg, ok := segs[segment]
				if !ok {
					seg = make(map[int]string)
				}
				seg[i] = val
				t.Logf("segment: %q, value: %q\n", segment, seg[i])
				segs[segment] = seg
			}
		}

		last = parser.d.state
		i = parser.d.hl7FieldIdx
	}

	require.NoError(t, parser.d.savedError)
	require.Contains(t, p.Keys(segs), "EVN")
	require.Equal(t, "A01", segs["EVN"][1])
}
