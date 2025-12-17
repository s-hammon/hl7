package v23

type EVN struct {
	EventTypeCode        string
	RecordedDt           string
	PlannedEventDateTime string
	EventReasonCode      string
	OperatorID           XCN
	EventOccurred        string
}

type PID struct {
	SetId                  string
	ExternalPatientId      CX
	InternalPatientId      CX
	AlternatePatientId     CX
	PatientName            XPN
	MotherMaidenName       XPN
	DOB                    string
	Sex                    string
	PatientAlias           XPN
	Race                   string
	PatientAddress         XAD
	CountyCode             string
	HomePhoneNumber        XTN
	WorkPhoneNumber        XTN
	PrimaryLanguage        CE
	MaritalStatus          string
	Religion               string
	PatientAccountNumber   CX
	SSN                    string
	DriversLicenseNumber   DLN
	MotherIdentifier       CX
	EthnicGroup            string
	BirthPlace             string
	MultipleBirthIndicator string
	BirthOrder             string
	Citizenship            string
	VeteranStatus          CE
	Nationality            CE
	PatientDeathDateTime   string
	PatientDeathIndicator  string
}

type PV1 struct {
	SetId                   string
	PatientClass            string
	AssignedPatientLocation PL
	AdmissionType           string
	PreadmitNumber          CX
	PriorPatientLocation    PL
	AttendingDoctor         XCN
	ReferringDoctor         XCN
	ConsultingDoctor        XCN
	HospitalService         string
	TemporaryLocation       PL
	PreadmitTestIndicator   string
	ReadmissionIndicator    string
	AdmitSource             string
	AmbulatoryStatus        string
	VipIndicator            string
	AdmittingDoctor         XCN
	PatientType             string
	VisitNumber             CX
	FinancialClass          FC
	ChargePriceIndicator    string
	CourtesyCode            string
	CreditRating            string
	ContractCode            string
	ContractEffectiveDate   string
	ContractAmount          string
	ContractPeriod          string
	InterestCode            string
	TransferBadDebtCode     string
	TransferBadDebtDate     string
	BadDebtAgencyCode       string
	BadDebtTransferAmount   string
	BadDebtRecoveryAmount   string
	DeleteAccountIndicator  string
	DeleteAccountDate       string
	DischargeDisposition    string
	DischargedToLocation    CM_DSL
	DietType                string
	ServicingFacility       string
	BedStatus               string
	AccountStatus           string
	PendingLocation         PL
	PriorTemporaryLocation  PL
	AdmitDateTime           string
	DischargeDateTime       string
	CurrentPatientBalance   string
	TotalCharges            string
	TotalAdjustments        string
	TotalPayments           string
	AlternateVisitId        CX
	VisitIndicator          string
	OtherHealthcareProvider XCN
}

type PV2 struct {
	PriorPendingLocation              PL
	AccomodationCode                  CE
	AdmitReason                       CE
	TransferReason                    CE
	PatientValuables                  string
	PatientValuablesLocation          string
	VisitUserCode                     string
	ExpectedAdmitDateTime             string
	ExpectedDischargeDateTime         string
	EstimatedLengthInpatientStay      string
	ActualLengthInpatientStay         string
	VisitDescription                  string
	ReferralSourceCode                XCN
	PreviousServiceDAte               string
	EmploymentIllnessRelatedIndicator string
	PurgeStatusCode                   string
	PurgeStatusDate                   string
	SpecialProgramCode                string
	RetentionIndicator                string
	ExpectedCountInsurancePlans       string
	VisitPublicityCode                string
	VisitProtectionIndicator          string
	ClinicOrganizationName            XON
	PatientStatusCode                 string
	VisitPriorityCode                 string
	PreviousTreatmentDAte             string
	ExpectedDischargeDisposition      string
	FileSignatureDate                 string
	FirstSimilarIllnessDate           string
	PatientChargeAdjustmentCode       string
	RecurringServiceCode              string
	BillingMediaCode                  string
	ExpectedSurgeryDateTime           string
	MilitaryPartnershipCode           string
	MilitaryNonAvailabilityCode       string
	NewbornBabyIndicator              string
	BabyDetainedIndicator             string
}

type PD1 struct {
	LivingDependency       string
	LivingArrangement      string
	PatientPrimaryFacility XON
	PatientPCPName         XCN
	StudentIndicator       string
	Handicap               string
	LivingWill             string
	OrganDonor             string
	SeparateBill           string
	DuplicatePatient       CX
	PublicityIndicator     CE
	ProtectionIndicator    string
}

type AL1 struct {
	SetId              string
	AllergyType        string
	AllergyCode        CE
	AllergySeverity    string
	AllergyReaction    string
	IdentificationDate string
}
