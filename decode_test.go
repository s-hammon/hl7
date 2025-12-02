package hl7

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeState_Unmarshal(t *testing.T) {
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
