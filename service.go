package main

import (
	"github.com/Financial-Times/tme-reader/tmereader"
	"github.com/pborman/uuid"
	"log"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type specialReportService interface {
	init() error
	getSpecialReportIds() ([]string)
	getSpecialReportsLinks() ([]specialReportLink, bool)
	getSpecialReportByUUID(uuid string) (specialReport, bool)
	checkConnectivity() error
}

type specialReportServiceImpl struct {
	repository         tmereader.Repository
	baseURL            string
	IdMap              map[string]string
	specialReportLinks []specialReportLink
	taxonomyName       string
	maxTmeRecords      int
}

func newSpecialReportService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (specialReportService, error) {

	s := &specialReportServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.init()
	if err != nil {
		return &specialReportServiceImpl{}, err
	}
	return s, nil
}

func (s *specialReportServiceImpl) init() error {
	s.IdMap = make(map[string]string)
	responseCount := 0
	log.Printf("Fetching special reports from TME\n")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Printf("Finished fetching special reports from TME\n")
			break
		}
		s.initSpecialReportsMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d special reports links\n", len(s.specialReportLinks))
	return nil
}

func (s *specialReportServiceImpl) getSpecialReportIds() []string {
	ids := make([]string, 0, len(s.IdMap))
	for id := range s.IdMap {
		ids = append(ids, id)
	}
	return ids
}

func (s *specialReportServiceImpl) getSpecialReportsLinks() ([]specialReportLink, bool) {
	if len(s.specialReportLinks) > 0 {
		return s.specialReportLinks, true
	}
	return s.specialReportLinks, false
}

func (s *specialReportServiceImpl) getSpecialReportByUUID(uuid string) (specialReport, bool) {
	rawId, found := s.IdMap[uuid]
	if !found {
		return specialReport{}, false
	}
	content, err := s.repository.GetTmeTermById(rawId)
	if err != nil {
		return specialReport{}, false
	}
	return transformSpecialReport(content.(term), s.taxonomyName), true
}

func (s *specialReportServiceImpl) checkConnectivity() error {
	// TODO: Can we just hit an endpoint to check if TME is available? Or do we need to make sure we get genre taxonmies back? Maybe a healthcheck or gtg endpoint?
	// TODO: Can we use a count from our responses while actually in use to trigger a healthcheck?
	//	_, err := s.repository.GetTmeTermsFromIndex(1)
	//	if err != nil {
	//		return err
	//	}
	return nil
}

func (s *specialReportServiceImpl) initSpecialReportsMap(terms []interface{}) {
	for _, iTerm := range terms {
		t := iTerm.(term)
		tmeIdentifier := buildTmeIdentifier(t.RawID, s.taxonomyName)
		uuid := uuid.NewMD5(uuid.UUID{}, []byte(tmeIdentifier)).String()
		s.IdMap[uuid] = t.RawID
		s.specialReportLinks = append(s.specialReportLinks, specialReportLink{APIURL: s.baseURL + uuid})
	}
}
