package hl7

import (
	"errors"
	"fmt"
)

type decodeState struct {
	data        []byte
	off         int // next read offset in data
	hl7FieldIdx int
	savedError  error

	fldDelim byte
	comDelim byte
	repDelim byte
	escDelim byte
	subDelim byte
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	d.hl7FieldIdx = 0
	d.savedError = nil

	return d
}

func (d *decodeState) unmarshal() error {
	if len(d.data) < 8 {
		return fmt.Errorf("not enough bytes in header: expecting at least 8, got %d", len(d.data))
	}

	d.fldDelim = d.data[3]
	d.comDelim = d.data[4]
	d.repDelim = d.data[5]
	d.escDelim = d.data[6]
	d.subDelim = d.data[7]

	d.off = 7
	d.hl7FieldIdx = 2
	return d.segment()
}

func (d *decodeState) segment() error {
	for d.off < len(d.data) {
		switch d.data[d.off] {
		case d.fldDelim:
			d.hl7FieldIdx++
		case byte('\r'):
			d.hl7FieldIdx = 0

			seg, err := d.peekSegName()
			if err != nil {
				return err
			}
			if seg == "MSH" || seg == "" {
				return nil
			}
		}

		d.off++
	}

	return nil
}

func (d *decodeState) peekSegName() (string, error) {
	i := d.off + 1
	if i >= len(d.data) {
		return "", nil
	}

	start := i
	for i < len(d.data) {
		if d.data[i] == d.fldDelim {
			return string(d.data[start:i]), nil
		}

		i++
	}

	return "", errors.New("unexpected end of data while scanning segment name")
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
