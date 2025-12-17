package v23

type MSH struct {
	FieldDelimiter           string
	EncodingCharacters       string
	SendingApplication       string
	SendingFacility          string
	ReceivingApplicatoin     string
	ReceivingFacility        string
	DateTime                 string
	Security                 string
	MessageType              CM_MSG
	ControlId                string
	ProcessingId             string
	VersionId                string
	SequenceNumber           string
	ContinuationPointer      string
	AcceptAcknowledgmentType string
	CountryCode              string
	CharacterSet             string
	PrincipalLanguage        string
}

type NTE struct {
	SetId           string
	SourceOfComment string
	Comment         string
}
