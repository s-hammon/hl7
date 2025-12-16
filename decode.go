package hl7

import (
	"fmt"
	"reflect"
	"strings"
)

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "hl7: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "hl7: Unmarshal(non-pointer " + e.Type.String() + ")"
	}

	return "hl7: Unmarshal(nil " + e.Type.String() + ")"
}

type decodeState struct {
	data       []byte
	off        int // next read offset in data
	prev       int // previous decoder state
	hl7Idx int // the current 1-based HL7 field index
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
	d.hl7Idx = 2
	return d
}

// after calling init
func (d *decodeState) unmarshal(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	rv = rv.Elem()

	var m map[string]any
	d.scanNext()
	if err := d.value(reflect.ValueOf(&m).Elem()); err != nil {
		return err
	}
	
	switch rv.Kind() {
	case reflect.Map:
		rv.Set(reflect.ValueOf(m))
	case reflect.Struct:
		return unmarshalStruct(rv, m)
	}

	return nil
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
			d.hl7Idx++
			return
		}
	}

	d.eof()
}

func (d *decodeState) readIndex() int {
	return d.off -1
}

func (d *decodeState) value(v reflect.Value) error {
	t := v.Type()

	switch v.Kind() {
	case reflect.Map:
		if v.IsNil() {
			v.Set(reflect.MakeMap(t))
		}
	}

	var (
		segmentName string = "MSH"
		fieldMap = make(map[int]any)
		inserted = false
	)

	fieldMap[1] = string(d.scan.fldDelim)
	fieldMap[2] = d.encodingChars()

	start := d.off

	for {
		d.scanValue()
		switch d.prev {
		case stateEOF, stateError:
			if !inserted {
				setSegmentValue(v, segmentName, fieldMap)
			}
			return d.savedError
		case stateFieldIdx:
			i := d.readIndex()
			raw := string(d.data[start:i])
			fieldMap[d.hl7Idx] =  d.buildFieldValue(raw)
		case stateEndSegment:
			i := d.readIndex()
			raw := string(d.data[start:i])
			fieldMap[d.hl7Idx] =  d.buildFieldValue(raw)

			if !inserted {
				setSegmentValue(v, segmentName, fieldMap)
				inserted = true
			}

			if d.off+3 > len(d.data) {
				return d.savedError
			}

			segmentName = string(d.data[d.off:d.off+3])
			d.scanN(3)

			fieldMap = make(map[int]any)
			inserted = false

			d.scanNext()
			d.hl7Idx = 0
		}

		start = d.off
	}
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

func (d *decodeState) buildFieldValue(raw string) any {
	// repeats
	if strings.IndexByte(raw, d.scan.repDelim) != -1 {
		parts := strings.Split(raw, string(d.scan.repDelim))
		out := make([]any, 0, len(parts))
		for _, p := range parts {
			out = append(out, d.buildFieldValue(p))
		}

		return out
	}
	// components
	if strings.IndexByte(raw, d.scan.comDelim) != -1 {
		parts := strings.Split(raw, string(d.scan.comDelim))
		m := make(map[int]any, len(parts))
		for i, p := range parts {
			m[i+1] = p
		}

		return m
	}

	// scalar
	return raw
}

func setSegmentValue(v reflect.Value, name string, fieldMap map[int]any) {
	key := reflect.ValueOf(name)
	existing := v.MapIndex(key)

	if !existing.IsValid() {
		v.SetMapIndex(key, reflect.ValueOf(fieldMap))
		return
	}

	existing = existing.Elem()

	switch existing.Kind() {
	case reflect.Map:
		slice := []map[int]any{
			existing.Interface().(map[int]any),
			fieldMap,
		}
		v.SetMapIndex(key, reflect.ValueOf(slice))
	case reflect.Slice:
		slice := existing.Interface().([]map[int]any)
		slice = append(slice, fieldMap)
		v.SetMapIndex(key, reflect.ValueOf(slice))
	}
}

func unmarshalStruct(dst reflect.Value, data map[string]any) error {
	t := dst.Type()
	for i := range t.NumField() {
		sf := t.Field(i)
		fv := dst.Field(i)

		segName := sf.Tag.Get("hl7")
		if segName == "" {
			segName = sf.Name
		}

		segData, ok := data[segName]
		if !ok {
			continue
		}

		assignSegment(fv, segData)
	}

	return nil
}

func assignSegment(dst reflect.Value, seg any) {
	switch v := seg.(type) {

	case map[int]any:
		assignSegmentStruct(dst, v)

	case []map[int]any:
		slice := reflect.MakeSlice(dst.Type(), 0, len(v))
		for _, m := range v {
			elem := reflect.New(dst.Type().Elem()).Elem()
			assignSegmentStruct(elem, m)
			slice = reflect.Append(slice, elem)
		}
		dst.Set(slice)
	}
}

func assignSegmentStruct(dst reflect.Value, fields map[int]any) {
	for i := range dst.NumField() {
		fv := dst.Field(i)
		hl7Idx := i + 1

		val, ok := fields[hl7Idx]
		if !ok {
			continue
		}

		assignValue(fv, val)
	}
}

func assignValue(dst reflect.Value, src any) {
	switch v := src.(type) {

	case string:
		if dst.Kind() == reflect.String {
			dst.SetString(v)
		}

	case []string:
		if dst.Kind() == reflect.Slice {
			slice := reflect.MakeSlice(dst.Type(), len(v), len(v))
			for i := range v {
				slice.Index(i).SetString(v[i])
			}
			dst.Set(slice)
		}

	case map[int]any:
		if dst.Kind() == reflect.Struct {
			assignSegmentStruct(dst, v)
		}

	case []any:
		if dst.Kind() == reflect.Slice {
			slice := reflect.MakeSlice(dst.Type(), 0, len(v))
			for _, elem := range v {
				e := reflect.New(dst.Type().Elem()).Elem()
				assignValue(e, elem)
				slice = reflect.Append(slice, e)
			}
			dst.Set(slice)
		}
	}
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
