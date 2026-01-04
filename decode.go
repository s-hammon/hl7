package hl7

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(data []byte, v any) error {
	var d decodeState
	d.init(data)
	if d.savedError != nil {
		return d.savedError
	}

	return d.unmarshal(v)
}

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
	hl7Idx     int // the current 1-based HL7 field index
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
	return d.off - 1
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
		fieldMap           = make(map[int]any)
		inserted           = false
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
			fieldMap[d.hl7Idx] = d.buildFieldValue(raw)
		case stateEndSegment:
			i := d.readIndex()
			raw := string(d.data[start:i])
			fieldMap[d.hl7Idx] = d.buildFieldValue(raw)

			if !inserted {
				setSegmentValue(v, segmentName, fieldMap)
				inserted = true
			}

			if d.off+3 > len(d.data) {
				return d.savedError
			}

			segmentName = string(d.data[d.off : d.off+3])
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
	if n := 1 + strings.Count(raw, string(d.scan.repDelim)); n > 1 {
		out := make([]any, 0, n)
		for p := range strings.SplitSeq(raw, string(d.scan.repDelim)) {
			out = append(out, d.buildFieldValue(p))
		}

		return out
	}
	// components
	if n := 1 + strings.Count(raw, string(d.scan.comDelim)); n > 1 {
		m := make(map[int]any, n)
		i := 1
		for p := range strings.SplitSeq(raw, string(d.scan.comDelim)) {
			m[i] = d.buildSubComponentValue(p)
			i++
		}

		return m
	}

	// subcomponents (or scaler if none exist)
	return d.buildSubComponentValue(raw)
}

func (d *decodeState) buildSubComponentValue(raw string) any {
	if n := 1 + strings.Count(raw, string(d.scan.subDelim)); n > 1 {
		m := make(map[int]any, n)
		i := 1
		for p := range strings.SplitSeq(raw, string(d.scan.subDelim)) {
			m[i] = p
			i++
		}

		return m
	}

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

		// skip unexported
		if !fv.CanSet() {
			continue
		}

		tag := parseTag(sf.Tag.Get("hl7"))

		// group field
		if tag.Options.Group() {
			switch fv.Kind() {
			case reflect.Slice, reflect.Pointer:
				buildGroupSlice(fv, data)
				continue
			}
		}

		// segment field
		segName := tag.Name
		if segName == "" {
			segName = sf.Name
		}

		if segData, ok := data[segName]; ok {
			assignSegment(fv, segData)
			continue
		}

		// nested struct
		switch fv.Kind() {
		case reflect.Struct:
			unmarshalStruct(fv, data)
		case reflect.Pointer:
			if fv.Type().Elem().Kind() != reflect.Struct {
				continue
			}
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			unmarshalStruct(fv.Elem(), data)
		}
	}

	return nil
}

func buildGroupSlice(dst reflect.Value, data map[string]any) {
	var groupType reflect.Type
	isSlice := dst.Kind() == reflect.Slice
	if isSlice {
		groupType = dst.Type().Elem()
	} else {
		groupType = dst.Type()
	}
	fields := groupFields(groupType)

	var anchor groupField
	for _, f := range fields {
		if f.Required {
			anchor = f
			break
		}
	}

	if anchor.Name == "" {
		return
	}

	segments := normalizeSegmentSlice(data[anchor.Name])
	count := len(segments)
	if !isSlice {
		if count > 1 {
			count = 1
		}
	}

	sliceType := reflect.SliceOf(groupType)
	out := reflect.MakeSlice(sliceType, count, count)
	for i := range count {
		var group reflect.Value
		if groupType.Kind() == reflect.Pointer {
			group = reflect.New(groupType.Elem())
			for _, f := range fields {
				assignGroupField(group.Elem().Field(f.Index), f, data, i)
			}
			unmarshalStruct(group.Elem(), data)
		} else {
			group = reflect.New(groupType).Elem()
			for _, f := range fields {
				assignGroupField(group.Field(f.Index), f, data, i)
			}
			unmarshalStruct(group, data)
		}
		out.Index(i).Set(group)
	}

	if isSlice {
		dst.Set(out)
	} else if count > 0 {
		dst.Set(out.Index(0))
	}
}

type groupField struct {
	Name     string
	Required bool
	Index    int
}

func groupFields(typ reflect.Type) []groupField {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	out := make([]groupField, 0, typ.NumField())

	for i := range typ.NumField() {
		sf := typ.Field(i)
		tag := parseTag(sf.Tag.Get("hl7"))

		if tag.Name == "" {
			continue
		}

		if sf.Type.Kind() == reflect.Slice {
			continue
		}

		out = append(out, groupField{
			Name:     tag.Name,
			Required: tag.Options.Required(),
			Index:    i,
		})
	}

	return out
}

func assignGroupField(dst reflect.Value, f groupField, data map[string]any, idx int) {
	seg, ok := data[f.Name]
	if !ok {
		return
	}

	segments := normalizeSegmentSlice(seg)
	if idx >= len(segments) {
		return
	}

	assignSegment(dst, segments[idx])
}

func normalizeSegmentSlice(seg any) []map[int]any {
	switch v := seg.(type) {
	default:
		return nil
	case map[int]any:
		return []map[int]any{v}
	case []map[int]any:
		return v
	}
}

func assignSegment(dst reflect.Value, seg any) {
	switch v := seg.(type) {
	case map[int]any:
		switch dst.Kind() {
		default:
			return
		case reflect.Struct:
			assignSegmentStruct(dst, v)
		case reflect.Pointer:
			if dst.IsNil() {
				dst.Set(reflect.New(dst.Type().Elem()))
			}
			assignSegmentStruct(dst.Elem(), v)
		case reflect.Slice:
			if dst.IsNil() {
				dst.Set(reflect.MakeSlice(dst.Type(), 0, 1))
			}
			elemType := dst.Type().Elem()
			var elem reflect.Value
			if elemType.Kind() == reflect.Pointer {
				elem = reflect.New(elemType.Elem())
				assignSegmentStruct(elem.Elem(), v)
			} else {
				elem = reflect.New(elemType).Elem()
				assignSegmentStruct(elem, v)
			}
			dst.Set(reflect.Append(dst, elem))
		}
	case []map[int]any:
		if dst.Kind() != reflect.Slice {
			return
		}

		n := len(v)
		slice := reflect.MakeSlice(dst.Type(), n, n)
		elemType := dst.Type().Elem()

		for i, m := range v {
			var elem reflect.Value
			if elemType.Kind() == reflect.Pointer {
				elem = reflect.New(elemType.Elem())
				assignSegmentStruct(elem.Elem(), m)
			} else {
				elem = reflect.New(elemType).Elem()
				assignSegmentStruct(elem, m)
			}
			slice.Index(i).Set(elem)
		}
		if dst.IsNil() {
			dst.Set(slice)
		} else {
			dst.Set(reflect.AppendSlice(dst, slice))
		}
	}
}

func assignSegmentStruct(dst reflect.Value, fields map[int]any) {
	t := dst.Type()

	idx := 1
	for i := range dst.NumField() {
		sf := t.Field(i)
		fv := dst.Field(i)

		if !fv.CanSet() {
			continue
		}

		var hl7Idx int
		if tag := sf.Tag.Get("hl7"); tag != "" {
			n, err := strconv.Atoi(tag)
			if err != nil {
				continue
			}

			hl7Idx = n
		} else {
			hl7Idx = idx
			idx++
		}

		val, ok := fields[hl7Idx]
		if !ok {
			continue
		}

		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			assignValue(fv.Elem(), val)
		} else {
			assignValue(fv, val)
		}
	}
}

func assignValue(dst reflect.Value, src any) {
	switch v := src.(type) {
	case string:
		switch dst.Kind() {
		case reflect.String:
			dst.SetString(v)
		case reflect.Struct:
			fields := map[int]any{1: v}
			assignSegmentStruct(dst, fields)
		case reflect.Pointer:
			if dst.IsNil() {
				dst.Set(reflect.New(dst.Type().Elem()))
			}
			fields := map[int]any{1: v}
			assignSegmentStruct(dst.Elem(), fields)
		}
	case map[int]any:
		switch dst.Kind() {
		case reflect.Struct:
			assignSegmentStruct(dst, v)
		case reflect.Pointer:
			if dst.IsNil() {
				dst.Set(reflect.New(dst.Type().Elem()))
			}
			assignSegmentStruct(dst.Elem(), v)
		}
	case []any:
		if dst.Kind() == reflect.Slice {
			n := len(v)
			slice := reflect.MakeSlice(dst.Type(), n, n)
			elemType := dst.Type().Elem()

			for i, elem := range v {
				var e reflect.Value
				if elemType.Kind() == reflect.Pointer {
					e = reflect.New(elemType.Elem())
				} else {
					e = reflect.New(elemType).Elem()
				}
				assignValue(e, elem)
				slice.Index(i).Set(e)
			}
			dst.Set(slice)
		}
	}
}
