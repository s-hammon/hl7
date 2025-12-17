package v23

type ORM struct {
	MSH          MSH
	PatientGroup PatientGroup
	OrderGroups  []OrderGroup `hl7:"group"`
}
