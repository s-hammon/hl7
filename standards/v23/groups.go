package v23

type PatientGroup struct {
	PID PID
	PD1 PD1
	NTE []NTE

	PatientVisitGroup PatientVisitGroup
	InsuranceGroup    []InsuranceGroup
	GT1               GT1
	AL1               []AL1
}

type PatientVisitGroup struct {
	PV1 PV1
	PV2 *PV2
}

type InsuranceGroup struct {
	IN1 IN1
	IN2 IN2
	IN3 IN3
}

type OrderGroup struct {
	ORC     ORC `hl7:"ORC,required"`
	Details OrderDetailGroup
}

type OrderDetailGroup struct {
	OBR OBR   `hl7:"OBR,required"`
	NTE []NTE `hl7:"NTE"`
	DG1 []DG1 `hl7:"DG1"`

	ObservationGroup []ObservationGroup `hl7:"group"`
}

type ObservationGroup struct {
	OBX OBX   `hl7:"OBX,required"`
	NTE []NTE `hl7:"NTE"`
}
