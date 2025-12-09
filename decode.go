package hl7

import (
	"fmt"
)

type decodeState struct {
	data       []byte
	off        int // next read offset in data
	prev       int
	scan       scanner
	savedError error
}

const (
	stateBegin int = iota
	stateHeaderSegment
	stateFieldIdx
	stateSegmentName
	stateEndSegment
	stateValue
	stateError
	stateEOF
)

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	d.prev = stateBegin

	if len(d.data) < 8 {
		d.savedError = fmt.Errorf("not enough bytes in header: expecting at least 8, got %d", len(d.data))
		return d
	}

	if string(d.data[:3]) != "MSH" {
		d.savedError = fmt.Errorf("expecting \"MSH\", got %q", string(d.data[:3]))
		return d
	}

	d.scan.fldDelim = d.data[3]
	d.scan.comDelim = d.data[4]
	d.scan.repDelim = d.data[5]
	d.scan.escDelim = d.data[6]
	d.scan.subDelim = d.data[7]

	d.off = 8
	return d
}

func (d *decodeState) encodingChars() string {
	chars := []byte{
		d.scan.comDelim,
		d.scan.repDelim,
		d.scan.escDelim,
		d.scan.subDelim,
	}
	return string(chars)
}

func (d *decodeState) scanNext() {
	if d.off < len(d.data) {
		d.prev = d.scan.state(d.data[d.off])
		d.off++
	} else {
		d.eof()
	}
}

func (d *decodeState) scanValue() {
	s, data, i := &d.scan, d.data, d.off
	for i < len(d.data) {
		current := s.state(data[i])
		i++
		if current != stateValue {
			d.prev = current
			d.off = i
			return
		}
	}

	d.eof()
}

func (d *decodeState) scanN(n int) {
	for range n {
		d.scanNext()
	}
}

func (d *decodeState) eof() {
	d.prev = stateEOF
	d.off = len(d.data) + 1
}

type ADT struct {
	MSH MSH `hl7:"MSH"`
	EVN EVN `hl7:"EVN"`
}

type MSH struct {
	FieldDelimiter     string
	EncodingCharacters string
	SendingApp         string
	SendingFac         string
	ReceivingApp       string
	ReceivingFac       string
	MessageDt          string
	Security           string
	MessageType        CM_MSG
	ControlId          string
	ProcessingId       string
	VersionId          string
}

type CM_MSG struct {
	Type         string
	TriggerEvent string
}

type EVN struct {
	EventTypeCode string
	RecordedDt    string
}

type ORM struct {
	MSH          MSH
	NTE          []NTE
	PatientGroup *PatientGroup
	OrderGroups  []OrderGroup
}

type PatientGroup struct {
	PID PID
	PD1 *PD1
	NTE []NTE

	PatientVisitGroup *PatientVisitGroup
	InsuranceGroup    []InsuranceGroup
	GT1               *GT1
	AL1               []AL1
}

type PatientVisitGroup struct {
	PV1 PV1
	PV2 *PV2
}

type InsuranceGroup struct {
	IN1 IN1
	IN2 *IN2
	IN3 *IN3
}

type OrderGroup struct {
	ORC              ORC
	OrderDetailGroup *OrderDetailGroup
}

type OrderDetailGroup struct {
	OBR OBR
	NTE []NTE
	DG1 []DG1

	ObservationGroup []ObservationGroup

	CTI *CTI
	BLG *BLG
}

type ObservationGroup struct {
	OBX OBX
	NTE []NTE
}

type GT1 struct{}
type AL1 struct{}
type PV1 struct{}
type PV2 struct{}
type IN1 struct{}
type IN2 struct{}
type IN3 struct{}
type ORC struct{}
type OBR struct{}
type DG1 struct{}
type BLG struct{}
type CTI struct{}
type OBX struct{}
type NTE struct{}
type PID struct{}
type PD1 struct{}
