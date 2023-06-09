package adp

import (
	"fmt"
	"strings"
)

type Worker struct {
	AssociateOID          string                  `json:"associateOID"`
	WorkerID              WorkerID                `json:"workerID"`
	Person                WorkerPersonInformation `json:"person"`
	WorkerDates           WorkerDates             `json:"workerDates"`
	WorkerStatus          WorkerStatus            `json:"workerStatus"`
	BusinessCommunication WorkerCommunication     `json:"businessCommunication"`
	WorkAssignments       []WorkerWorkAssignment  `json:"workAssignments"`
}

func (w *Worker) GetAssociateOID() string {
	return w.AssociateOID
}

func (w *Worker) GetFirstName() string {
	return w.Person.LegalName.GivenName
}

func (w *Worker) GetMiddleName() string {
	return w.Person.LegalName.MiddleName
}

func (w *Worker) GetLastName() string {
	return w.Person.LegalName.FamilyName
}

func (w *Worker) GetBirthDate() WorkerDate {
	return ToWorkerDate(w.Person.BirthDate)
}

func (w *Worker) GetGender() string {
	gender := w.Person.GenderCode.LongName
	if gender == "" {
		gender = w.Person.GenderCode.ShortName
	}
	return gender
}

func (w *Worker) IsDisabled() bool {
	return w.Person.DisabledIndicator
}

func (w *Worker) GetEthnicity() string {
	ethnicity := w.Person.EthnicityCode.LongName
	if ethnicity == "" {
		ethnicity = w.Person.EthnicityCode.ShortName
	}
	return ethnicity
}

func (w *Worker) ListGovernmentIDs() []WorkerGovernmentID {
	return w.Person.GovernmentIDs
}

func (w *Worker) GetAddress() WorkerLegalAddress {
	return w.Person.LegalAddress
}

func (w *Worker) GetJobTitle() *string {
	assignment := w.GetPrimaryWorkAssignment()
	if assignment == nil {
		return nil
	}
	jobTitle := assignment.GetJobTitle()
	return &jobTitle
}

func (w *Worker) ListWorkAssignments() []WorkerWorkAssignment {
	return w.WorkAssignments
}

func (w *Worker) GetPrimaryWorkAssignment() *WorkerWorkAssignment {
	for _, assignment := range w.WorkAssignments {
		if assignment.PrimaryIndicator == true {
			return &assignment
		}
	}
	return nil
}

func (w *Worker) IsActive() bool {
	return w.WorkerStatus.StatusCode["codeValue"] == "Active"
}

func (w *Worker) GetBusinessEmails() []string {
	emails := w.BusinessCommunication.Emails
	emailAddresses := []string{}
	for _, email := range emails {
		emailAddresses = append(emailAddresses, email.EmailUri)
	}

	return emailAddresses
}

func (w *Worker) GetPersonnalEmails() []string {
	emails := w.Person.Communication.Emails
	emailAddresses := []string{}
	for _, email := range emails {
		emailAddresses = append(emailAddresses, email.EmailUri)
	}

	return emailAddresses
}

type WorkerDate struct {
	Year  string
	Month string
	Day   string
}

func (w WorkerDate) String() string {
	return fmt.Sprintf("%s-%s-%s", w.Year, w.Month, w.Day)
}

func ToWorkerDate(date string) WorkerDate {
	splitDate := strings.Split(date, "-")
	return WorkerDate{
		Year:  splitDate[0],
		Month: splitDate[1],
		Day:   splitDate[2],
	}
}

type WorkerID struct {
	IDValue string `json:"idValue"`
}

type WorkerPersonInformation struct {
	BirthDate                  string                         `json:"birthDate"`
	GenderCode                 WorkerGender                   `json:"genderCode"`
	MartialStatus              WorkerMaritalStatusCode        `json:"maritalStatus"`
	SocialInsurancePrograms    []WorkerSocialInsuranceProgram `json:"socialInsurancePrograms"`
	TobaccoUserIndicator       bool                           `json:"tobacooUserIndicator"`
	DisabledIndicator          bool                           `json:"disabledIndicator"`
	EthnicityCode              WorkerEthnicity                `json:"ethnicityCode"`
	MilitaryStatusCode         map[string]string              `json:"militaryStatusCode"`
	MilitaryDischargeDate      string                         `json:"militaryDischargeDate"`
	MilitaryClassifiationCodes []map[string]string            `json:"militaryClassificationCodes"`
	GovernmentIDs              []WorkerGovernmentID           `json:"governmentIDs"`
	LegalName                  WorkerLegalName                `json:"legalName"`
	LegalAddress               WorkerLegalAddress             `json:"legalAddress"`
	Communication              WorkerCommunication            `json:"communication"`
}

type WorkerEthnicity struct {
	CodeValue string `json:"codeValue"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

type WorkerGender struct {
	CodeValue string `json:"codeValue"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

type WorkerMaritalStatusCode struct {
	EffectiveDate string `json:"effectiveDate"`
	CodeValue     string `json:"codeValue"`
	ShortName     string `json:"shortName"`
}

type WorkerSocialInsuranceProgram struct {
	NameCode         map[string]string `json:"codeValue"`
	CoveredIndicator bool              `json:"coveredIndicator"`
}

type WorkerGovernmentID struct {
	ItemID      string            `json:"itemID"`
	IDValue     string            `json:"idValue"`
	NameCode    map[string]string `json:"nameCode"`
	CountryCode string            `json:"countryCode"`
}

type WorkerLegalName struct {
	GenerationAffixCode    map[string]string `json:"generationAffixCode"`
	QualificationAffixCode map[string]string `json:"qualificationAffixCode"`
	GivenName              string            `json:"givenName"`
	MiddleName             string            `json:"middleName"`
	FamilyName             string            `json:"familyName1"`
	FormattedName          string            `json:"formattedName"`
}

type WorkerLegalAddress struct {
	NameCode            map[string]string `json:"nameCode"`
	LineOne             string            `json:"lineOne"`
	LineTwo             string            `json:"lineTwo"`
	CityName            string            `json:"cityName"`
	CountrySubDivision1 string            `json:"countrySubdivision1"`
	CountrySubDivision2 string            `json:"countrySubdivision2"`
	CountryCode         string            `json:"countryCode"`
	PostalCode          string            `json:"postalCode"`
}

type WorkerLegalAddressCountrySubDivision struct {
	SubDivisiontype string `json:"subdivisionType"`
	CodeValue       string `json:"codeValue"`
	ShortName       string `json:"shortName"`
}

type WorkerCommunication struct {
	Emails []WorkerCommunicationEmails `json:"emails"`
}

type WorkerCommunicationEmails struct {
	NameCode map[string]string `json:"nameCode"`
	EmailUri string            `json:"emailUri"`
}

type WorkerDates struct {
	OriginalHireDate string `json:"originalHireDate"`
}

type WorkerStatus struct {
	StatusCode map[string]string `json:"statusCode"`
}

type WorkerWorkAssignment struct {
	ItemID                      string                         `json:"itemID"`
	PrimaryIndicator            bool                           `json:"primaryIndicator"`
	HireDate                    string                         `json:"hireDate"`
	ActualStartDate             string                         `json:"actualStartDate"`
	AssignmentStatus            WorkerWorkAssignmentStatus     `json:"assignmentStatus"`
	JobCode                     map[string]string              `json:"jobCode"`
	JobTitle                    string                         `json:"jobTitle"`
	OccupationalClassifications []map[string]map[string]string `json:"occupationalClassifications"`
	IndustryClassifications     []map[string]map[string]string `json:"industryClassifications"`
	WageLawCoverage             map[string]map[string]string   `json:"wageLawCoverage"`
	PositionId                  string                         `json:"positionID"`
	AssignedWorkLocations       []WorkerAssignedWorkLocation   `json:"assignedWorkLocation"`
	ReportsTo                   []WorkerReportsTo              `json:"reportsTo"`
}

func (w *WorkerWorkAssignment) GetJobTitle() string {
	return w.JobTitle
}

func (w *WorkerWorkAssignment) ListReportsToAssociateOID() []string {
	results := []string{}
	for _, reportsTo := range w.ReportsTo {
		results = append(results, reportsTo.AssociateOID)
	}
	return results
}

func (w *WorkerWorkAssignment) GetHireDate() WorkerDate {
	return ToWorkerDate(w.HireDate)
}

func (w *WorkerWorkAssignment) GetActualStartDate() WorkerDate {
	return ToWorkerDate(w.ActualStartDate)
}

type WorkerWorkAssignmentStatus struct {
	StatusCode    map[string]string `json:"statusCode"`
	ReasonCode    map[string]string `json:"reasonCode"`
	EffectiveDate string            `json:"effectiveDate"`
}

type WorkerAssignedWorkLocation struct {
	Addresss WorkerAssignedWorkLocationAddress `json:"address"`
}

type WorkerAssignedWorkLocationAddress struct {
	NameCode                 map[string]string `json:"nameCode"`
	LineOne                  string            `json:"lineOne"`
	CityName                 string            `json:"cityName"`
	CountrySubdivisionLevel1 map[string]string `json:"countrySubdivision1"`
	CountryCode              string            `json:"countryCode"`
}

type WorkerReportsTo struct {
	PositionID          string                  `json:"positionID"`
	AssociateOID        string                  `json:"associateOID"`
	WorkerID            WorkerReportsToWorkerID `json:"workerID"`
	ReportsToWorkerName map[string]string       `json:"reportsToWorkerName"`
}

type WorkerReportsToWorkerID struct {
	IDValue    string            `json:"idValue"`
	SchemaCode map[string]string `json:"workerID"`
}
