package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name          string
		term          term
		specialReport specialReport
	}{
		{"Transform term to specialReport", term{
			CanonicalName: "Test Special Report",
			RawID:         "Nstein_GL_AFTM_GL_164835"},
			specialReport{
				UUID:      "adb4f804-c3b6-3eca-8708-5edeec653a27",
				PrefLabel: "Test Special Report",
				AlternativeIdentifiers: alternativeIdentifiers{
					TME:   []string{"TnN0ZWluX0dMX0FGVE1fR0xfMTY0ODM1-U2VjdGlvbnM="},
					Uuids: []string{"adb4f804-c3b6-3eca-8708-5edeec653a27"},
				},
				Type: "SpecialReport"}},
	}

	for _, test := range tests {
		expectedSection := transformSpecialReport(test.term, "Sections")
		assert.Equal(test.specialReport, expectedSection, fmt.Sprintf("%s: Expected specialReport incorrect", test.name))
	}

}
