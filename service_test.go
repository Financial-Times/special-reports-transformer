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
		name      string
		baseURL   string
		terms     []term
		specialReports []specialReportLink
		found     bool
		err       error
	}{
		{"Success", "localhost:8080/transformers/special-reports/",
			[]term{term{CanonicalName: "Business Guide to Manchester 2012", RawID: "8f57aae4-d6f1-4322-88b0-9d9f176ffd8c"}},
			[]specialReportLink{specialReportLink{APIURL: "localhost:8080/transformers/special-reports/a39edaf1-534d-33b8-9e6b-d80137e75ef8"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/special-reports/", []term{}, []specialReportLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newSpecialReportService(&repo, test.baseURL, "SpecialReports", 10000)
		actualSpecialReports, found := service.getSpecialReports()
		assert.Equal(test.specialReports, actualSpecialReports, fmt.Sprintf("%s: Expected special report link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetSpecialReportByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		terms    []term
		uuid     string
		specialReport specialReport
		found    bool
		err      error
	}{
		{"Success", []term{term{CanonicalName: "Business Guide to Manchester 2012", RawID: "8f57aae4-d6f1-4322-88b0-9d9f176ffd8c"}},
			"a39edaf1-534d-33b8-9e6b-d80137e75ef8", specialReport{UUID: "a39edaf1-534d-33b8-9e6b-d80137e75ef8", CanonicalName: "Business Guide to Manchester 2012", TmeIdentifier: "OGY1N2FhZTQtZDZmMS00MzIyLTg4YjAtOWQ5ZjE3NmZmZDhj-U3BlY2lhbFJlcG9ydHM=", Type: "Special Report"}, true, nil},
		{"Not found", []term{term{CanonicalName: "Business Guide to Manchester 2012", RawID: "8f57aae4-d6f1-4322-88b0-9d9f176ffd8c"}},
			"some uuid", specialReport{}, false, nil},
		{"Error on init", []term{}, "some uuid", specialReport{}, false, errors.New("Error getting taxonomy")},
	}
	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newSpecialReportService(&repo, "", "SpecialReports", 10000)
		actualSpecialReport, found := service.getSpecialReportByUUID(test.uuid)
		assert.Equal(test.specialReport, actualSpecialReport, fmt.Sprintf("%s: Expected special report incorrect", test.name))
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
