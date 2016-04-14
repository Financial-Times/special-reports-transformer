package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		term     term
		specialReport specialReport
	}{
		{"Transform term to special report", term{CanonicalName: "Business Guide to Manchester 2012", RawID: "8f57aae4-d6f1-4322-88b0-9d9f176ffd8c"}, specialReport{UUID: "a39edaf1-534d-33b8-9e6b-d80137e75ef8", CanonicalName: "Business Guide to Manchester 2012", TmeIdentifier: "OGY1N2FhZTQtZDZmMS00MzIyLTg4YjAtOWQ5ZjE3NmZmZDhj-U3BlY2lhbFJlcG9ydHM=", Type: "Special Report"}},
	}

	for _, test := range tests {
		expectedSpecialReport := transformSpecialReport(test.term, "SpecialReports")

		assert.Equal(test.specialReport, expectedSpecialReport, fmt.Sprintf("%s: Expected special report incorrect", test.name))
	}

}
