package hl7

import (
	"testing"

	v23 "github.com/s-hammon/hl7/proto/standards/v23"
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
	msg := []byte("MSH|^~\\&|SendingApp|SendingFac|ReceivingApp|ReceivingFac|20250101000000||ORM^O01|123456|P|2.3|4232072\rPID|1||V12345||DOE^JANE^A||19700101|F|||123 MAIN ST^ANYWHERE^TX^76543^USA||(123)456-7890\rPV1||E|Acme ER^AER^^AR||||123456^Smith^John^J^^^M.D.\rORC|XO|00112233|30504059||CM||^^^20250101080000||20250101100000|^Decrad^Support^^^^System.||123456^Smith^John^J^^^M.D.|LTERRAD1^LT ER RAD1\rOBR|1|00112233|30504059|CXR^Chest 1 View|Y^N||20250101000000")

	var m v23.ORM_O01
	d := newState(msg)
	err := d.unmarshal(&m)
	require.NoError(t, err)
	require.Equal(t, "SendingApp", m.MSH.SendingApplication)
	require.Equal(t, "SendingFac", m.MSH.SendingFacility)
	require.Equal(t, "SendingFac", m.MSH.SendingFacility)
	require.Equal(t, "ORM", m.MSH.MessageType.Type)
	require.Equal(t, "O01", m.MSH.MessageType.TriggerEvent)
	require.Equal(t, "V12345", m.PatientGroup.PID.InternalPatientId.Id)
	require.Equal(t, "XO", m.OrderGroups[0].ORC.OrderControl)
	require.Equal(t, "30504059", m.OrderGroups[0].ORC.FillerOrderNumber)
	require.Equal(t, "20250101080000", m.OrderGroups[0].ORC.QuantityTiming.StartDateTime)
}

func TestUnmarshal_ORM(t *testing.T) {
	msg := []byte("MSH|^~\\&|SendingApp|SendingFac|ReceivingApp|ReceivingFac|20250101000000||ORM^O01|123456|P|2.3|4232072\rPID|1||V12345||DOE^JANE^A||19700101|F|||123 MAIN ST^ANYWHERE^TX^76543^USA||(123)456-7890\rPV1||E|Acme ER^AER^^AR||||123456^Smith^John^J^^^M.D.\rORC|XO|00112233|30504059||CM||^^^20250101080000||20250101100000|^Decrad^Support^^^^System.||123456^Smith^John^J^^^M.D.|LTERRAD1^LT ER RAD1\rOBR|1|00112233|30504059|CXR^Chest 1 View|Y^N||20250101000000\r")

	var (
		m   v23.ORM_O01
		err error
	)

	err = Unmarshal(msg, m)
	require.Error(t, err)

	err = Unmarshal(msg, &m)
	require.NoError(t, err)
}

func TestUnmarshal_ORU(t *testing.T) {
	msg := []byte("MSH|^~\\&|PSOne|BMCNE|STRIC|STRIC|20250404152739||ORU^R01|6767683|P|2.3|29069747\rPID|||002207830||SMITH^JINKLEHEIMER^JOHN JACOB||19840526|M|||123 MAIN STR^^ANYWHERE^TX^12345^USA||(999)999-9999|(999)999-9999\rPV1||O|Boutique Mammography Center at^BMCNE^^BMCNE^^^^^ACME MAMMOGRAPHY CENTER||||440854^DOE^JANE^^^^M.D.^^NPI&1234567890|||||||||||O||||||||||||||||||||||||||20250404000100\rORC|RE||29737914||||20250404152445^20250404152445^20250404152535\rOBR|1|12|29737914|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD |||20250404152445||||||No Known Allergies|||440854^DOE^JANE^^^^M.D.^^NPI&1234567890||V00384534|D01620528|  -  ,   -  ,   -  |STRICAH051|20250404152535|STRICAH051|MG|F||^^^20250404115000^20250404115900||||SCR|620863&Farkas&Julie&M&&&M.D.^^20250404152535\rOBX|1|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||MAMMOGRAM DIGITAL SCREENING BILATERAL W/CAD AND DBT||||||F|||20250404152535\rOBX|2|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|3|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||DATE:  4/4/2025||||||F|||20250404152535\rOBX|4|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|5|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||HISTORY:  Screening||||||F|||20250404152535\rOBX|6|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||  ||||||F|||20250404152535\rOBX|7|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||TECHNIQUE:  Bilateral full field digital screening mammography and bilateral||||||F|||20250404152535\rOBX|8|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||digital breast tomosynthesis were performed and interpreted in conjunction with||||||F|||20250404152535\rOBX|9|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||computer-aided detection. CC and MLO views were obtained of the breast(s) with||||||F|||20250404152535\rOBX|10|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||additional views as required. ||||||F|||20250404152535\rOBX|11|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|12|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||COMPARISON: Prior mammograms dating back to 11/16/2019 ||||||F|||20250404152535\rOBX|13|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|14|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||FINDINGS:||||||F|||20250404152535\rOBX|15|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|16|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||RIGHT: No suspicious mass, suspicious architectural distortion, or suspicious||||||F|||20250404152535\rOBX|17|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||microcalcifications. ||||||F|||20250404152535\rOBX|18|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|19|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||LEFT: No suspicious mass, suspicious architectural distortion, or suspicious||||||F|||20250404152535\rOBX|20|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||microcalcifications. ||||||F|||20250404152535\rOBX|21|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|22|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||OTHER: None.||||||F|||20250404152535\rOBX|23|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|24|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||IMPRESSION:||||||F|||20250404152535\rOBX|25|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|26|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||No suspicious findings. Unless otherwise indicated, continue annual screening||||||F|||20250404152535\rOBX|27|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||mammogram.||||||F|||20250404152535\rOBX|28|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|29|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||BIRADS Category 1 - Negative ||||||F|||20250404152535\rOBX|30|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|31|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||Your patient is being notified by mail of the results.||||||F|||20250404152535\rOBX|32|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|33|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||Breast Density:  The breasts are heterogeneously dense, which may obscure small||||||F|||20250404152535\rOBX|34|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||masses (Type C)||||||F|||20250404152535\rOBX|35|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|36|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||MFC:  1NC||||||F|||20250404152535\rOBX|37|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535\rOBX|38|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||Signed on 4/4/2025 3:25 PM by Julie M Farkas, M.D.||||||F|||20250404152535\rOBX|39|FT|MAMSTOM2^Mammogram Digital Screening Bilateral w/CAD ||||||||F|||20250404152535")

	var (
		m   v23.ORU_R01
		err error
	)

	err = Unmarshal(msg, &m)
	require.NoError(t, err)
	t.Log(m)
	require.Len(t, m.Results, 1)
	require.Len(t, m.Results[0].Order[0].Observation, 39)
	require.Equal(t, "MAMMOGRAM DIGITAL SCREENING BILATERAL W/CAD AND DBT", m.Results[0].Order[0].Observation[0].OBX.ObservationValue)
	require.Equal(t, "Signed on 4/4/2025 3:25 PM by Julie M Farkas, M.D.", m.Results[0].Order[0].Observation[37].OBX.ObservationValue)
	require.Equal(t, "", m.Results[0].Order[0].Observation[38].OBX.ObservationValue)
	require.Equal(t, "O", m.Results[0].Visit.PV1.PatientClass)
	require.Equal(t, "JINKLEHEIMER", m.Results[0].PID.PatientName.GivenName)
	require.Equal(t, "19840526", m.Results[0].PID.Dob)
}

func TestUnmarshal_ORU_MultipleOrders(t *testing.T) {
	msg := []byte("MSH|^~\\&|PSOne|METHNE||MHS|20251216002851||ORU^R01|7853152|P|2.3|4232074\rPID|||V00272475||BANANA^ANNA^BANNA||19801006|F|||123 MAIN ST^^ANYWHERE^TX^76543^USA||(123)456-7890|(098)765-4321||||V468357251\rPV1||E|NEMH ER FT 757 5009^VFT^^METHNE^^^^^NEMH ER FT 757 5009 VFT|||||||||||||||E|V468357251\rORC|RE||30507023||||20251216002425^20251216002425^20251216002644\rOBR|1|002353470|30507023|UPELNOB^US Pelvis Non-OB|||20251216002425||||||Lower abd pain, r ovarian cyst    DX:  ABD PAIN    Comments:  #V468357251|||123456^Smigh^John^A^^^P.A.||N00069961||, , |RAD-DOCTOR|20251216002644|RAD-DOCTOR|US|F||^^^20251215234700^20251215234700||||Lower abd pain, r ovarian cyst|999696&Graham&Joshua&J&&&M.D.^^20251216002644\rORC|CN||30507022||||20251216002425^20251216002425^20251216002644\rOBR|2|002353469|30507022|UPELDOP^US Doppler Pelvis|||20251216002425||||||Lower abd pain, r ovarian cyst    DX:  ABD PAIN    Comments:  #V468357251|||123456^Smigh^John^A^^^P.A.||N00069961||, , |Methodist Hospital Northeast|20251216002644|RAD-DOCTOR|US|F||^^^20251215234700^20251215234700||||Lower abd pain, r ovarian cyst|999696&Graham&Joshua&J&&&M.D.^^20251216002644\rOBX|1|FT|UPELNOB^US Pelvis Non-OB||ULTRASOUND PELVIS ||||||F|||20251216002644\rOBX|2|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|3|FT|UPELNOB^US Pelvis Non-OB||DATE:  12/15/2025||||||F|||20251216002644\rOBX|4|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|5|FT|UPELNOB^US Pelvis Non-OB||HISTORY:   Lower abd pain, r ovarian cyst    ||||||F|||20251216002644\rOBX|6|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|7|FT|UPELNOB^US Pelvis Non-OB||TECHNIQUE: Ultrasound of the pelvis performed per the routine protocol using a||||||F|||20251216002644\rOBX|8|FT|UPELNOB^US Pelvis Non-OB||transabdominal and transvaginal probe.||||||F|||20251216002644\rOBX|9|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|10|FT|UPELNOB^US Pelvis Non-OB||COMPARISON: CT dated 12/15/2025||||||F|||20251216002644\rOBX|11|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|12|FT|UPELNOB^US Pelvis Non-OB||FINDINGS:||||||F|||20251216002644\rOBX|13|FT|UPELNOB^US Pelvis Non-OB|| ||||||F|||20251216002644\rOBX|14|FT|UPELNOB^US Pelvis Non-OB||Uterus: Not seen ||||||F|||20251216002644\rOBX|15|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|16|FT|UPELNOB^US Pelvis Non-OB||Right Ovary: 4.4 x 4.6 x 3.4 cm||||||F|||20251216002644\rOBX|17|FT|UPELNOB^US Pelvis Non-OB||Cysts: 2.8 x 3.1 x 2.7 cm||||||F|||20251216002644\rOBX|18|FT|UPELNOB^US Pelvis Non-OB||Mass: None||||||F|||20251216002644\rOBX|19|FT|UPELNOB^US Pelvis Non-OB||Doppler examination: No evidence for ovarian torsion. Normal spectral Doppler||||||F|||20251216002644\rOBX|20|FT|UPELNOB^US Pelvis Non-OB||waveforms with pulsatile arterial inflow and aphasic venous outflow.||||||F|||20251216002644\rOBX|21|FT|UPELNOB^US Pelvis Non-OB|| ||||||F|||20251216002644\rOBX|22|FT|UPELNOB^US Pelvis Non-OB||Left Ovary: Not seen||||||F|||20251216002644\rOBX|23|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|24|FT|UPELNOB^US Pelvis Non-OB||Free fluid: None||||||F|||20251216002644\rOBX|25|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|26|FT|UPELNOB^US Pelvis Non-OB||IMPRESSION: ||||||F|||20251216002644\rOBX|27|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|28|FT|UPELNOB^US Pelvis Non-OB||Right ovarian cyst, corresponds to CT finding.||||||F|||20251216002644\rOBX|29|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|30|FT|UPELNOB^US Pelvis Non-OB||Signed on 12/16/2025 12:26 AM by Joshua J Graham, M.D.||||||F|||20251216002644\rOBX|31|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\r")

	var (
		m   v23.ORU_R01
		err error
	)

	err = Unmarshal(msg, &m)
	require.NoError(t, err)
	require.Len(t, m.Results, 1)
	require.Len(t, m.Results[0].Order, 2)
	// require.Len(t, m.OrderGroup, 1)
	t.Log(m.Results[0].Order[0].OBR)
	require.Equal(t, "30507023", m.Results[0].Order[0].ORC.FillerOrderNumber)
	require.Equal(t, "30507023", m.Results[0].Order[0].OBR.FillerOrderNumber)
	require.Equal(t, "", m.Results[0].Order[0].OBR.Priority)
	require.Equal(t, "999696", m.Results[0].Order[0].OBR.PrincipalResultInterpreter.Name.IdNumber)
	require.Equal(t, "30507022", m.Results[0].Order[1].ORC.FillerOrderNumber)
	require.Equal(t, "30507022", m.Results[0].Order[1].OBR.FillerOrderNumber)
	require.Equal(t, "", m.Results[0].Order[1].OBR.Priority)
	require.Len(t, m.Results[0].Order[0].Observation, 31)
	// TODO: fix this so that we only assign the OBX segments once
	require.Len(t, m.Results[0].Order[1].Observation, 31)
}

func BenchmarkUnmarshal_ORU(b *testing.B) {
	msg := []byte("MSH|^~\\&|PSOne|METHNE||MHS|20251216002851||ORU^R01|7853152|P|2.3|4232074\rPID|||V00272475||BANANA^ANNA^BANNA||19801006|F|||123 MAIN ST^^ANYWHERE^TX^76543^USA||(123)456-7890|(098)765-4321||||V468357251\rPV1||E|NEMH ER FT 757 5009^VFT^^METHNE^^^^^NEMH ER FT 757 5009 VFT|||||||||||||||E|V468357251\rORC|RE||30507023||||20251216002425^20251216002425^20251216002644\rOBR|1|002353470|30507023|UPELNOB^US Pelvis Non-OB|||20251216002425||||||Lower abd pain, r ovarian cyst    DX:  ABD PAIN    Comments:  #V468357251|||123456^Smigh^John^A^^^P.A.||N00069961||, , |RAD-DOCTOR|20251216002644|RAD-DOCTOR|US|F||^^^20251215234700^20251215234700||||Lower abd pain, r ovarian cyst|999696&Graham&Joshua&J&&&M.D.^^20251216002644\rORC|CN||30507022||||20251216002425^20251216002425^20251216002644\rOBR|2|002353469|30507022|UPELDOP^US Doppler Pelvis|||20251216002425||||||Lower abd pain, r ovarian cyst    DX:  ABD PAIN    Comments:  #V468357251|||123456^Smigh^John^A^^^P.A.||N00069961||, , |Methodist Hospital Northeast|20251216002644|RAD-DOCTOR|US|F||^^^20251215234700^20251215234700||||Lower abd pain, r ovarian cyst|999696&Graham&Joshua&J&&&M.D.^^20251216002644\rOBX|1|FT|UPELNOB^US Pelvis Non-OB||ULTRASOUND PELVIS ||||||F|||20251216002644\rOBX|2|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|3|FT|UPELNOB^US Pelvis Non-OB||DATE:  12/15/2025||||||F|||20251216002644\rOBX|4|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|5|FT|UPELNOB^US Pelvis Non-OB||HISTORY:   Lower abd pain, r ovarian cyst    ||||||F|||20251216002644\rOBX|6|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|7|FT|UPELNOB^US Pelvis Non-OB||TECHNIQUE: Ultrasound of the pelvis performed per the routine protocol using a||||||F|||20251216002644\rOBX|8|FT|UPELNOB^US Pelvis Non-OB||transabdominal and transvaginal probe.||||||F|||20251216002644\rOBX|9|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|10|FT|UPELNOB^US Pelvis Non-OB||COMPARISON: CT dated 12/15/2025||||||F|||20251216002644\rOBX|11|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|12|FT|UPELNOB^US Pelvis Non-OB||FINDINGS:||||||F|||20251216002644\rOBX|13|FT|UPELNOB^US Pelvis Non-OB|| ||||||F|||20251216002644\rOBX|14|FT|UPELNOB^US Pelvis Non-OB||Uterus: Not seen ||||||F|||20251216002644\rOBX|15|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|16|FT|UPELNOB^US Pelvis Non-OB||Right Ovary: 4.4 x 4.6 x 3.4 cm||||||F|||20251216002644\rOBX|17|FT|UPELNOB^US Pelvis Non-OB||Cysts: 2.8 x 3.1 x 2.7 cm||||||F|||20251216002644\rOBX|18|FT|UPELNOB^US Pelvis Non-OB||Mass: None||||||F|||20251216002644\rOBX|19|FT|UPELNOB^US Pelvis Non-OB||Doppler examination: No evidence for ovarian torsion. Normal spectral Doppler||||||F|||20251216002644\rOBX|20|FT|UPELNOB^US Pelvis Non-OB||waveforms with pulsatile arterial inflow and aphasic venous outflow.||||||F|||20251216002644\rOBX|21|FT|UPELNOB^US Pelvis Non-OB|| ||||||F|||20251216002644\rOBX|22|FT|UPELNOB^US Pelvis Non-OB||Left Ovary: Not seen||||||F|||20251216002644\rOBX|23|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|24|FT|UPELNOB^US Pelvis Non-OB||Free fluid: None||||||F|||20251216002644\rOBX|25|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|26|FT|UPELNOB^US Pelvis Non-OB||IMPRESSION: ||||||F|||20251216002644\rOBX|27|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|28|FT|UPELNOB^US Pelvis Non-OB||Right ovarian cyst, corresponds to CT finding.||||||F|||20251216002644\rOBX|29|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\rOBX|30|FT|UPELNOB^US Pelvis Non-OB||Signed on 12/16/2025 12:26 AM by Joshua J Graham, M.D.||||||F|||20251216002644\rOBX|31|FT|UPELNOB^US Pelvis Non-OB||||||||F|||20251216002644\r")

	var m v23.ORU_R01

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = Unmarshal(msg, &m)
	}
}

func TestUnmarshal_AL1(t *testing.T) {
	msg := []byte("MSH|^~\\&|ITS|WOH|METHWO|METHWO|202512220000||ORM^O01|19738904|P|2.4\rPID|1|L1-B20250520173705489|Q267415244^^^METHWO^^METHWO||DOE^JOHN||19700601|M|||123 MAIN ST^^ANYWHERE^TX^76543||123-456-7890|||M|BAP|A26740438416\rPV1|1|I|A.ICU^A.IC07A^A^METHWO^^^^^A.ICU A.IC07A|EM|||^House^Gregory^^^^DO|^Referred^Self||ICU||||PR|||DNE7747^House^Gregory^^^^DO|I|A26740438416|08|||||||||||||||||||COCWH||ADM|||202512201650\rAL1|1|DA|F006004444^ranitidine^^From Zantac^^allergy.id|MO|Rash|20251220\rORC|NW|A000000078463A|A000000078463A||SC|N|^^^202512220504^^R||202512220000|||OJL9891^Farkas^Julie^^^^MD|MWORM1|210-690-7400|||A.ICU\rOBR|1|A000000078463A|A000000078463A|MHXRCXR1V^XR chest 1V^MWORM1|R|202512220504|202512220504||||||pneumonia|||005845^Farkas^Julie^^^^MD|210-690-7400^^PH^^^210^690-7400||Q267415244|A26740438416|Methodist Hospital Westover Hills|||XR|||1^^^202512220504^^R|UNKNOWN^MISSING^NUMBER~SELF^Referred^Self|||pneumonia|||^^^^MWORM1|||||Julie  Farkas  MD  -  210-690-7400\rOBX|1|TX|ORDERPTTYPE||I\rOBX|1|CE|MHXRCXR1V^XR chest 1V||H")

	var (
		m   v23.ORM_O01
		err error
	)

	err = Unmarshal(msg, &m)
	require.NoError(t, err)
	t.Log(m)
	require.Len(t, m.PatientGroup.AL1, 1)
	require.Equal(t, "1", m.PatientGroup.AL1[0].SetId)
	require.Equal(t, "ranitidine", m.PatientGroup.AL1[0].AllergyCode.Text)
}
