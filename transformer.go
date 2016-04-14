package main

import (
	"encoding/base64"
	"encoding/xml"
	"github.com/pborman/uuid"
)

func transformSpecialReport(tmeTerm term, taxonomyName string) specialReport {
	tmeIdentifier := buildTmeIdentifier(tmeTerm.RawID, taxonomyName)

	return specialReport{
		UUID:          uuid.NewMD5(uuid.UUID{}, []byte(tmeIdentifier)).String(),
		CanonicalName: tmeTerm.CanonicalName,
		TmeIdentifier: tmeIdentifier,
		Type:          "Special Report",
	}
}

func buildTmeIdentifier(rawId string, tmeTermTaxonomyName string) string {
	id := base64.StdEncoding.EncodeToString([]byte(rawId))
	taxonomyName := base64.StdEncoding.EncodeToString([]byte(tmeTermTaxonomyName))
	return id + "-" + taxonomyName
}

type specialReportsTransformer struct {
}

func (*specialReportsTransformer) UnMarshallTaxonomy(contents []byte) ([]interface{}, error) {
	taxonomy := taxonomy{}
	err := xml.Unmarshal(contents, &taxonomy)
	if err != nil {
		return nil, err
	}
	var interfaces []interface{} = make([]interface{}, len(taxonomy.Terms))
	for i, d := range taxonomy.Terms {
		interfaces[i] = d
	}
	return interfaces, nil
}

func (*specialReportsTransformer) UnMarshallTerm(content []byte) (interface{}, error) {
	tmeTerm := term{}
	err := xml.Unmarshal(content, &tmeTerm)
	if err != nil {
		return term{}, err
	}
	return tmeTerm, nil
}
