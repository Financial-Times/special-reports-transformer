package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSpecialReports(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name           string
		baseURL        string
		terms          []term
		specialreports []specialReportLink
		found          bool
		err            error
	}{
		{"Success", "localhost:8080/transformers/specialreports/",
			[]term{term{CanonicalName: "Z_Archive", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "mNGQ2MWQ0NDMtMDc5Mi00NWExLTlkMGQtNWZhZjk0NGExOWU2-Z2VucVz"}},
			[]specialReportLink{specialReportLink{APIURL: "localhost:8080/transformers/specialreports/20ddda23-a1bb-3530-88aa-60232583895a"},
				specialReportLink{APIURL: "localhost:8080/transformers/specialreports/cfd7a2d5-bc8f-3585-b98a-db69f7b8cfea"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/specialreports/", []term{}, []specialReportLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newSpecialReportService(&repo, test.baseURL, "Sections", 10000)
		expectedSections, found := service.getSpecialReports()
		assert.Equal(test.specialreports, expectedSections, fmt.Sprintf("%s: Expected SpecialReports link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetSectionByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name          string
		terms         []term
		uuid          string
		specialreport specialReport
		found         bool
		err           error
	}{
		{"Success", []term{term{CanonicalName: "SpecialReport1", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "SpecialReport2", RawID: "TkdRMk1XUTBORE10TURjNU1pMDBOV0V4TFRsa01HUXROV1poWmprME5HRXhPV1UyLVoyVnVjbVZ6-U2VjdGlvbnM=]"}},
			"ccd5cc74-1f1b-3ac6-a563-e36dff51926c", getDummySpecialReport("ccd5cc74-1f1b-3ac6-a563-e36dff51926c", "SpecialReport1", "YjgzMzc1NTktYWMwOC0zNDA0LTkwMjUtYmFkNTFlYmUyZmM3-U3BlY2lhbFJlcG9ydHM="), true, nil},
		{"Not found", []term{term{CanonicalName: "SpecialReport3", RawID: "845dc7d7-ae89-4fed-a819-9edcbb3fe507"}, term{CanonicalName: "Feature", RawID: "NGQ2MWdefsdfsfcmVz"}},
			"some uuid", specialReport{}, false, nil},
		{"Error on init", []term{}, "some uuid", specialReport{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newSpecialReportService(&repo, "", "SpecialReports", 10000)
		expectedSpecialReport, found := service.getSpecialReportByUUID(test.uuid)
		assert.Equal(test.specialreport, expectedSpecialReport, fmt.Sprintf("%s: Expected SpecialReports incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	terms []term
	err   error
}

func (d *dummyRepo) GetTmeTermsFromIndex(startRecord int) ([]interface{}, error) {
	if startRecord > 0 {
		return nil, d.err
	}
	var interfaces []interface{} = make([]interface{}, len(d.terms))
	for i, data := range d.terms {
		interfaces[i] = data
	}
	return interfaces, d.err
}
func (d *dummyRepo) GetTmeTermById(uuid string) (interface{}, error) {
	return d.terms[0], d.err
}

func getDummySpecialReport(uuid string, prefLabel string, tmeId string) specialReport {
	return specialReport{
		UUID:      uuid,
		PrefLabel: prefLabel,
		Type:      "SpecialReport",
		AlternativeIdentifiers: alternativeIdentifiers{TME: []string{tmeId}, Uuids: []string{uuid}}}
}
