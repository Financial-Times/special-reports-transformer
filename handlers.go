package main

import (
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const space = " "

type specialReportsHandler struct {
	service specialReportService
}

type specialReportId struct {
	id string
}

func newSpecialReportsHandler(service specialReportService) specialReportsHandler {
	return specialReportsHandler{service: service}
}

func (h *specialReportsHandler) getCount(writer http.ResponseWriter, req *http.Request) {
	ids := h.service.getSpecialReportIds()
	writer.Write([]byte(strconv.Itoa(len(ids))))
}

func (h *specialReportsHandler) getIds(writer http.ResponseWriter, req *http.Request) {
	ids := h.service.getSpecialReportIds()
	writer.Header().Add("Content-Type", "text/plain")
	if len(ids) == 0 {
		writer.WriteHeader(http.StatusOK)
		return
	}
	enc := json.NewEncoder(writer)
	for _, id := range ids {
		rId := specialReportId{id: id}
		err := enc.Encode(rId)
		if err != nil {
			log.Warnf("Couldn't encode to HTTP response special report with uuid=%s %v\n", id, err)
			continue
		}
		_, err = writer.Write([]byte(space))
		if err != nil {
			log.Warnf("Couldn't write ' ' to HTTP response for special report with uuid=%s %v\n", id, err)
			continue
		}
	}
}

func (h *specialReportsHandler) reload(writer http.ResponseWriter, req *http.Request) {
	err := h.service.init()
	if err != nil {
		log.Warnf("Problem reloading terms from TME: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *specialReportsHandler) getSpecialReports(writer http.ResponseWriter, req *http.Request) {
	obj, found := h.service.getSpecialReportsLinks()
	writeJSONResponse(obj, found, writer)
}

func (h *specialReportsHandler) getSpecialReportByUUID(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	obj, found := h.service.getSpecialReportByUUID(uuid)
	writeJSONResponse(obj, found, writer)
}

func writeJSONResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v\n", err)
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", errorMsg))
}

// HealthCheck does something
func (h *specialReportsHandler) HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for the special report data from TME",
		Name:             "Check connectivity to TME",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/specialReports-transfomer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to TME to be able to supply special reports",
		Checker:          h.checker,
	}
}

// Checker does more stuff
func (h *specialReportsHandler) checker() (string, error) {
	err := h.service.checkConnectivity()
	if err == nil {
		return "Connectivity to TME is ok", err
	}
	return "Error connecting to TME", err
}

//GoodToGo returns a 503 if the healthcheck fails - suitable for use from varnish to check availability of a node
func (h *specialReportsHandler) GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := h.checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}
