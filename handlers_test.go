package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getSpecialReportResponse = "[{\"apiUrl\":\"http://localhost:8080/transformers/special-reports/bba39990-c78d-3629-ae83-808c333c6dbc\"}]\n"
const getSpecialReportByUUIDResponse = "{\"uuid\":\"bba39990-c78d-3629-ae83-808c333c6dbc\",\"canonicalName\":\"Metals Markets\",\"tmeIdentifier\":\"MTE3-U3ViamVjdHM=\",\"type\":\"Special Report\"}\n"

func TestHandlers(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		req          *http.Request
		dummyService specialReportService
		statusCode   int
		contentType  string // Contents of the Content-Type header
		body         string
	}{
		{"Success - get special report by uuid", newRequest("GET", fmt.Sprintf("/transformers/special-reports/%s", testUUID)), &dummyService{found: true, specialReports: []specialReport{specialReport{UUID: testUUID, CanonicalName: "Metals Markets", TmeIdentifier: "MTE3-U3ViamVjdHM=", Type: "Special Report"}}}, http.StatusOK, "application/json", getSpecialReportByUUIDResponse},
		{"Not found - get special report by uuid", newRequest("GET", fmt.Sprintf("/transformers/special-reports/%s", testUUID)), &dummyService{found: false, specialReports: []specialReport{specialReport{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get special reports", newRequest("GET", "/transformers/special-reports"), &dummyService{found: true, specialReports: []specialReport{specialReport{UUID: testUUID}}}, http.StatusOK, "application/json", getSpecialReportResponse},
		{"Not found - get special reports", newRequest("GET", "/transformers/special-reports"), &dummyService{found: false, specialReports: []specialReport{}}, http.StatusNotFound, "application/json", ""},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.True(test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.Equal(test.body, rec.Body.String(), fmt.Sprintf("%s: Wrong body", test.name))
	}
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func router(s specialReportService) *mux.Router {
	m := mux.NewRouter()
	h := newSpecialReportsHandler(s)
	m.HandleFunc("/transformers/special-reports", h.getSpecialReports).Methods("GET")
	m.HandleFunc("/transformers/special-reports/{uuid}", h.getSpecialReportByUUID).Methods("GET")
	return m
}

type dummyService struct {
	found     bool
	specialReports []specialReport
}

func (s *dummyService) getSpecialReports() ([]specialReportLink, bool) {
	var links []specialReportLink
	for _, sub := range s.specialReports {
		links = append(links, specialReportLink{APIURL: "http://localhost:8080/transformers/special-reports/" + sub.UUID})
	}
	return links, s.found
}

func (s *dummyService) getSpecialReportByUUID(uuid string) (specialReport, bool) {
	return s.specialReports[0], s.found
}
