package main

import (
	"github.com/Financial-Times/tme-reader/tmereader"
	"github.com/pborman/uuid"
	"log"
)

type specialReportService interface {
	reload() error
	getCount() int
	getIds() []string
	getLinks() ([]specialReportLink, bool)
	getSpecialReport(uuid string) (specialReport, bool)
}

type specialReportServiceImpl struct {
	repository         tmereader.Repository
	baseURL            string
	IDMap              map[string]string
	specialReportLinks []specialReportLink
	taxonomyName       string
	maxTmeRecords      int
}

func newSpecialReportService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (specialReportService, error) {
	s := &specialReportServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.reload()
	if err != nil {
		return &specialReportServiceImpl{}, err
	}
	return s, nil
}

func (s *specialReportServiceImpl) reload() error {
	s.IDMap = make(map[string]string)
	var links []specialReportLink
	s.specialReportLinks = links
	responseCount := 0
	log.Println("Fetching special reports from TME")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}
		if len(terms) < 1 {
			log.Println("Finished fetching special reports from TME")
			break
		}
		s.initSpecialReportsMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d special reports links\n", len(s.specialReportLinks))
	return nil
}

func (s *specialReportServiceImpl) getCount() int {
	return len(s.IDMap)
}

func (s *specialReportServiceImpl) getIds() []string {
	ids := make([]string, 0, len(s.IDMap))
	for id := range s.IDMap {
		ids = append(ids, id)
	}
	return ids
}

func (s *specialReportServiceImpl) getLinks() ([]specialReportLink, bool) {
	if len(s.specialReportLinks) > 0 {
		return s.specialReportLinks, true
	}
	return s.specialReportLinks, false
}

func (s *specialReportServiceImpl) getSpecialReport(uuid string) (specialReport, bool) {
	rawID, found := s.IDMap[uuid]
	if !found {
		return specialReport{}, false
	}
	content, err := s.repository.GetTmeTermById(rawID)
	if err != nil {
		return specialReport{}, false
	}
	return transformSpecialReport(content.(term), s.taxonomyName), true
}

func (s *specialReportServiceImpl) initSpecialReportsMap(terms []interface{}) {
	for _, iTerm := range terms {
		t := iTerm.(term)
		tmeIdentifier := buildTmeIdentifier(t.RawID, s.taxonomyName)
		uuid := uuid.NewMD5(uuid.UUID{}, []byte(tmeIdentifier)).String()
		s.IDMap[uuid] = t.RawID
		s.specialReportLinks = append(s.specialReportLinks, specialReportLink{APIURL: s.baseURL + uuid})
	}
}
