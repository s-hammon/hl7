package v23

type ORM_O01 struct {
	MSH          MSH
	NTE          NTE
	PatientGroup PatientGroup
	OrderGroups  OrderGroup `hl7:"group"`
}

type ORU_R01 struct {
	MSH     MSH
	Results []ResultGroup `hl7:"group"`
	DCS     DSC
}
