package hl7

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newState(msg []byte) *decodeState {
	d := &decodeState{}
	d.init(msg)
	return d
}

func TestDecodeState_UnmarshalError(t *testing.T) {
	var a map[string]any
	d := newState([]byte("MSH|^~\\&|\r"))
	err := d.unmarshal(a)
	require.Error(t, err)

	var invalidErr *InvalidUnmarshalError
	require.ErrorAs(t, err, &invalidErr)

	var b *map[string]any = nil
	d = newState([]byte("MSH|^~\\&|\r"))
	err = d.unmarshal(b)
	require.Error(t, err)
	require.ErrorAs(t, err, &invalidErr)

	d = newState([]byte("ABC|123"))
	var c map[string]any
	err = d.unmarshal(&c)
	require.Error(t, err)
	require.NotErrorAs(t, err, &invalidErr)
}

func TestDecodeState_UnmarshalMap(t *testing.T) {
	msg := []byte("MSH|^~\\&|SendingApp|SendingFac|ORU^R01\rPID|1|123|Doe, Jane~Smith, John\rOBX|1|FT|CXR^Chest 1 View\rOBX|2|FT|CXR^Chest 1 View\r")

	var m map[string]any
	d := newState(msg)

	want := map[string]any{
		"MSH": map[int]any{
			1: "|",
			2: "^~\\&",
			3: "SendingApp",
			4: "SendingFac",
			5: map[int]any{
				1: "ORU",
				2: "R01",
			},
		},
		"PID": map[int]any{
			1: "1",
			2: "123",
			3: []any{
				"Doe, Jane",
				"Smith, John",
			},
		},
		"OBX": []map[int]any{
			{
				1: "1",
				2: "FT",
				3: map[int]any{
					1: "CXR",
					2: "Chest 1 View",
				},
			},
			{
				1: "2",
				2: "FT",
				3: map[int]any{
					1: "CXR",
					2: "Chest 1 View",
				},
			},
		},
	}

	startIdx := d.hl7Idx
	err := d.unmarshal(&m)
	require.NoError(t, err)
	require.GreaterOrEqual(t, d.hl7Idx, startIdx)
	require.Len(t, m, 3)
	require.Contains(t, m, "MSH")
	require.Contains(t, m, "PID")
	require.Contains(t, m, "OBX")
	require.Equal(t, want, m)
}

func TestDecodeState_UnmarshalORM(t *testing.T) {
	msg := []byte("MSH|^~\\&|SendingApp|SendingFac|ReceivingApp|ReceivingFac|20250101000000||ORM^O01|123456|P|2.3|4232072\rPID|1||V12345||DOE^JANE^A||19700101|F|||123 MAIN ST^ANYWHERE^TX^76543^USA||(123)456-7890\rPV1||E|Acme ER^AER^^AR||||123456^Smith^John^J^^^M.D.\rORC|XO|00112233|30504059||CM||^^^20250101080000||20250101100000|^Decrad^Support^^^^System.||123456^Smith^John^J^^^M.D.|LTERRAD1^LT ER RAD1\rOBR|1|00112233|30504059|CXR^Chest 1 View|Y^N||20250101000000\r")

	var m ORM
	d := newState(msg)
	err := d.unmarshal(&m)
	require.NoError(t, err)
	require.Equal(t, "SendingApp", m.MSH.SendingApp)
	require.Equal(t, "SendingFac", m.MSH.SendingFac)
	require.Equal(t, "SendingFac", m.MSH.SendingFac)
	t.Log(m)
	require.Equal(t, "ORM", m.MSH.MessageType.Type)
	require.Equal(t, "O01", m.MSH.MessageType.TriggerEvent)
	require.Equal(t, "O01", m.MSH.MessageType.TriggerEvent)
	require.Equal(t, "V12345", m.PatientGroup.PID.ExternalPatientId.Id)
	require.Len(t, m.OrderGroups, 1)
	require.Equal(t, "XO", m.OrderGroups[0].ORC.OrderControl)
	require.Equal(t, "30504059", m.OrderGroups[0].ORC.FillerOrderNumber)
	require.Equal(t, "20250101080000", m.OrderGroups[0].ORC.QuantityTiming.StartDt)
}
