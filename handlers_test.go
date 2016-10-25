package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getSpecialReportsResponse = `[{"apiUrl":"http://localhost:8080/transformers/specialreports/bba39990-c78d-3629-ae83-808c333c6dbc"}]`
const getSpecialReportByUUIDResponse = `{"uuid":"bba39990-c78d-3629-ae83-808c333c6dbc","alternativeIdentifiers":{"TME":["MTE3-U3ViamVjdHM="],"uuids":["bba39990-c78d-3629-ae83-808c333c6dbc"]},"prefLabel":"Global Special Reports","type":"SpecialReport"}`

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
		{"Success - get special report by uuid", newRequest("GET", fmt.Sprintf("/transformers/specialreports/%s", testUUID)), &dummyService{found: true, specialreports: []specialReport{getDummySpecialReport(testUUID, "Global Special Reports", "MTE3-U3ViamVjdHM=")}}, http.StatusOK, "application/json", getSpecialReportByUUIDResponse},
		{"Not found - get special report by uuid", newRequest("GET", fmt.Sprintf("/transformers/specialreports/%s", testUUID)), &dummyService{found: false, specialreports: []specialReport{specialReport{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get special reports", newRequest("GET", "/transformers/specialreports"), &dummyService{found: true, specialreports: []specialReport{specialReport{UUID: testUUID}}}, http.StatusOK, "application/json", getSpecialReportsResponse},
		{"Not found - get special reports", newRequest("GET", "/transformers/specialreports"), &dummyService{found: false, specialreports: []specialReport{}}, http.StatusNotFound, "application/json", ""},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.True(test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.Equal(strings.TrimSpace(test.body), strings.TrimSpace(rec.Body.String()), fmt.Sprintf("%s: Wrong body", test.name))
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
	m.HandleFunc("/transformers/specialreports", h.getSpecialReports).Methods("GET")
	m.HandleFunc("/transformers/specialreports/{uuid}", h.getSpecialReportByUUID).Methods("GET")
	return m
}

type dummyService struct {
	found          bool
	specialreports []specialReport
}

func (s *dummyService) init() error {
	return nil
}

func (s *dummyService) getSpecialReportIds() []string {
	var ids []string
	for _, sub := range s.specialreports {
		ids = append(ids, sub.UUID)
	}
	return ids
}

func (s *dummyService) getSpecialReportsLinks() ([]specialReportLink, bool) {
	var specialReportLinks []specialReportLink
	for _, sub := range s.specialreports {
		specialReportLinks = append(specialReportLinks, specialReportLink{APIURL: "http://localhost:8080/transformers/specialreports/" + sub.UUID})
	}
	return specialReportLinks, s.found
}

func (s *dummyService) getSpecialReportByUUID(uuid string) (specialReport, bool) {
	return s.specialreports[0], s.found
}

func (s *dummyService) checkConnectivity() error {
	return nil
}
