package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
	"github.com/sethgrid/pester"
	"net"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFormatter(new(log.JSONFormatter))
}

func main() {
	app := cli.App("special-reports-transformer", "A RESTful API for transforming TME SpecialReports to UP json")
	username := app.String(cli.StringOpt{
		Name:   "tme-username",
		Value:  "",
		Desc:   "TME username used for http basic authentication",
		EnvVar: "TME_USERNAME",
	})
	password := app.String(cli.StringOpt{
		Name:   "tme-password",
		Value:  "",
		Desc:   "TME password used for http basic authentication",
		EnvVar: "TME_PASSWORD",
	})
	token := app.String(cli.StringOpt{
		Name:   "token",
		Value:  "",
		Desc:   "Token to be used for accessig TME",
		EnvVar: "TOKEN",
	})
	baseURL := app.String(cli.StringOpt{
		Name:   "base-url",
		Value:  "http://localhost:8080/transformers/special-reports/",
		Desc:   "Base url",
		EnvVar: "BASE_URL",
	})
	tmeBaseURL := app.String(cli.StringOpt{
		Name:   "tme-base-url",
		Value:  "https://tme.ft.com",
		Desc:   "TME base url",
		EnvVar: "TME_BASE_URL",
	})
	port := app.Int(cli.IntOpt{
		Name:   "port",
		Value:  8080,
		Desc:   "Port to listen on",
		EnvVar: "PORT",
	})
	maxRecords := app.Int(cli.IntOpt{
		Name:   "maxRecords",
		Value:  int(10000),
		Desc:   "Maximum records to be queried to TME",
		EnvVar: "MAX_RECORDS",
	})
	slices := app.Int(cli.IntOpt{
		Name:   "slices",
		Value:  int(10),
		Desc:   "Number of requests to be executed in parallel to TME",
		EnvVar: "SLICES",
	})

	tmeTaxonomyName := "SpecialReports"

	app.Action = func() {
		client := getResilientClient()

		mf := new(specialReportsTransformer)
		s, err := newSpecialReportService(tmereader.NewTmeRepository(client, *tmeBaseURL, *username, *password, *token, *maxRecords, *slices, tmeTaxonomyName, &tmereader.AuthorityFiles{}, mf), *baseURL, tmeTaxonomyName, *maxRecords)
		if err != nil {
			log.Errorf("Error while creating SpecialReportService: [%v]", err.Error())
		}

		h := newSpecialReportsHandler(s)
		m := mux.NewRouter()

		m.HandleFunc(status.PingPath, status.PingHandler)
		m.HandleFunc(status.PingPathDW, status.PingHandler)
		m.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
		m.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)
		m.HandleFunc("/__health", v1a.Handler("Special Reports Transformer Healthchecks", "Checks for accessing TME", h.HealthCheck()))
		m.HandleFunc("/__gtg", h.GoodToGo)

		m.HandleFunc("/transformers/special-reports", h.getSpecialReports).Methods("GET")
		m.HandleFunc("/transformers/special-reports/__count", h.getCount).Methods("GET")
		m.HandleFunc("/transformers/special-reports/__ids", h.getIds).Methods("GET")
		m.HandleFunc("/transformers/special-reports/__reload", h.reload).Methods("POST")
		m.HandleFunc("/transformers/special-reports/{uuid}", h.getSpecialReport).Methods("GET")

		http.Handle("/", m)

		log.Printf("listening on %d", *port)
		err = http.ListenAndServe(fmt.Sprintf(":%d", *port),
			httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry,
				httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), m)))
		if err != nil {
			log.Errorf("Can't start up HTTP listener: %v", err)
		}
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("Cannot start app: %v", err)
	}
}

func getResilientClient() *pester.Client {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 32,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	c := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
	client := pester.NewExtendedClient(c)
	client.Backoff = pester.ExponentialBackoff
	client.MaxRetries = 5
	client.Concurrency = 1

	return client
}
