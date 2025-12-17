package v23

type OBX struct {
	SetId                        string
	ValueType                    string
	ObservationIdentifier        CE
	ObservationSubId             string
	ObservationValue             string
	Units                        CE
	ReferencesRange              string
	AbnormalFlags                string
	Probability                  string
	AbnormalTestNature           string
	ResultStatus                 string
	LastDateObservedNormalValues string
	UserDefinedAccessChecks      string
	ObservationDateTime          string
	ProducerId                   CE
	ResponsibleObserver          XCN
	ObservationMethod            CE
}
