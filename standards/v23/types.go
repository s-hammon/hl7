package v23

type CM_MSG struct {
	Type         string
	TriggerEvent string
}

type XCN struct {
	IdNumber             string
	FamilyName           string
	GivenName            string
	MiddleName           string
	Suffix               string
	Prefix               string
	Degree               string
	SourceTable          string
	AssigningAuthority   string
	NameTypeCode         string
	IdentifierCheckDigit string
	CheckDigitSchemeCode string
	IdentifierTypeCode   string
	AssigningFacility    string
}

type CX struct {
	Id                       string
	CheckId                  string
	CheckDigitIdentifierCode string
	AssigningAuthority       string
	IdentifierTypeCode       string
	AssigningFacility        string
}

type XPN struct {
	FamilyName             string
	GivenName              string
	MiddleName             string
	Suffix                 string
	Prefix                 string
	Degree                 string
	NameTypeCode           string
	NameRepresentationCode string
}

type XAD struct {
	StreetAddress              string
	OtherDesignation           string
	City                       string
	State                      string
	Zip                        string
	Country                    string
	AddressType                string
	OtherGeographicDesignation string
	CountyCode                 string
	CensusTract                string
}

type XTN struct {
	Number                   string
	TelecommunicationUseCode string
	EquipmentType            string
	EmailAddress             string
	CountryCode              string
	AreaCode                 string
	PhoneNumber              string
	Extension                string
	AnyText                  string
}

type CE struct {
	Identifier            string
	Text                  string
	CodingSystem          string
	AlternateIdentifier   string
	AlternateText         string
	AlternateCodingSystem string
}

type DLN struct {
	LicenseNumber  string
	IssuingState   string
	ExpirationDate string
}

type PL struct {
	PointOfCare         string
	Room                string
	Bed                 string
	Facility            string
	LocationStatus      string
	PersonLocationType  string
	Building            string
	Floor               string
	LocationDescription string
}

type FC struct {
	FinancialClass string
	EffectiveDate  string
}

type CM_DSL struct {
	Location      string
	EffectiveDate string
}

type XON struct {
	OrganizationName     string
	TypeCode             string
	IdNumber             string
	CheckDigit           string
	CheckDigitSchemeCode string // HL7 0061
	AssigningAuthority   string
	IdentifierTypeCode   string
	AssigningFacility    string
}

type CM_POR struct {
	PlacerOrderNumber string
	FillerOrderNumber string
}

type CQ struct {
	Quantity string
	Units    CE
}

type CP struct {
	Price      string
	PriceType  string
	FromValue  string
	ToValue    string
	RangeUnits CE
	RangeType  string
}

type JCC struct {
	JobCode  string
	JobClass string
}

type CM_AUI struct {
	AuthorizationNumber string
	Date                string
	Source              string
}

type CM_PLT struct {
	RoomType       string
	AmountType     string
	CoverageAmount string
}

type CM_DDE struct {
	DelayDays string
	Amount    string
	DayCount  string
}

type CM_VAL struct {
	Type   string
	Amount string
}

type CM_PCR struct {
	PatientType string
	Required    string
	Window      string
}

type CM_SPE struct {
	Name                         CE
	Additives                    string
	Freetext                     string
	BodySite                     CE
	SiteModifier                 CE
	CollectionMethodModifierCode CE
}

type CM_CHP struct {
	DollarAmount MO
	ChargeCOde   CE
}

type MO struct {
	Quantity     string
	Denomination string
}

type CM_PRE struct {
	ObservationIdentifier CE
	SubId                 string
	ObservationResult     string
}

type CM_OBS struct {
	Name                CN
	StartDateTime       string
	EndDateTime         string
	PointOfCare         string
	Room                string
	Bed                 string
	Facility            HD
	LocationStatus      string
	PatientLocationType string
	Building            string
	Floor               string
}

type CN struct {
	IdNumber           string
	FamilyName         string
	GivenName          string
	MiddleName         string
	Suffix             string
	Prefix             string
	Degree             string
	SourceTable        string
	AssigningAuthority HD
}

type HD struct {
	NamespaceId     string
	UniversalId     string
	UniversalIdType string
}
