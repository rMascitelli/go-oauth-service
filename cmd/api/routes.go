package main

import (
	"fmt"
	proclient "github.com/rMascitelli/go-prometheus-metrics-helper/client"
	"log"
	"net/http"
	"time"
)

//OAuth Specification is described in these RFC Articles:
//	https://www.rfc-editor.org/rfc/rfc6749
//	https://www.rfc-editor.org/rfc/rfc7662

func NewRouter(port int, postgres PostgresConnector) Router {
	if port <= 0 {
		log.Fatalf("Cannot create server at port %d\n", port)
	}

	return Router{
		port:          port,
		handler:       NewHandler(postgres),
		metricsClient: proclient.NewPrometheusClient(),
	}
}

func (r *Router) AddMetrics() {
	r.metricsClient.AddNewCounter(NUM_REGISTRY, "Number of registry requests")
	r.metricsClient.AddNewCounter(NUM_LOGIN, "Number of login requests")
	r.metricsClient.AddNewCounter(NUM_INTROSPECT, "Number of introspect requests")
	r.metricsClient.AddNewGauge(ELAPSED_REGISTRY_MS, "Time spent on last registry request")
	r.metricsClient.AddNewGauge(ELAPSED_LOGIN_MS, "Time spent on last login request")
	r.metricsClient.AddNewGauge(ELAPSED_INTROSPECT_MS, "Time spent on last introspect request")
}

func (r *Router) StartRouter() {
	log.Printf("Serving at port %d...\n", r.port)
	r.AddMetrics()
	http.HandleFunc("/register", r.Register)
	http.HandleFunc("/login", r.Login)
	http.HandleFunc("/introspect", r.Introspect)
	http.ListenAndServe(fmt.Sprintf(":%d", r.port), nil)
}

func (r *Router) Login(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	err, loginResponse := r.handler.HandleUserLogin(req)
	if err != nil {
		log.Printf("Error occured while handling User Login, err:\n  %v\n", err)
		writeJSONResponse(w, 400, err.Error())
	} else {
		log.Printf("Succesfully authenticated")
		writeJSONResponse(w, 200, loginResponse)
	}
	elapsed_s := float64(time.Since(start).Microseconds()) / 1000
	log.Printf("Elapsed login serve time - %f\n", elapsed_s)
	r.metricsClient.IncrementCounter(NUM_LOGIN, SERVICENAME)
	r.metricsClient.SetGaugeVal(ELAPSED_LOGIN_MS, SERVICENAME, elapsed_s)
}

func (r *Router) Register(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	if err := r.handler.HandleRegistry(req); err != nil {
		log.Printf("Error occured while handling Registry, err:\n  %v\n", err)
		writeJSONResponse(w, 400, err.Error())
	} else {
		log.Printf("Succesfully registered")
		writeJSONResponse(w, 200, "Success!")
	}
	elapsed_s := float64(time.Since(start).Microseconds()) / 1000
	log.Printf("Elapsed register serve time - %.3f\n", elapsed_s)
	r.metricsClient.IncrementCounter(NUM_REGISTRY, SERVICENAME)
	r.metricsClient.SetGaugeVal(ELAPSED_REGISTRY_MS, SERVICENAME, elapsed_s)
}

func (r *Router) Introspect(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	err, introspectResponse := r.handler.HandleIntrospect(req)
	if err != nil {
		log.Printf("Error occured while handling Introspect, err:\n  %v\n", err)
		writeJSONResponse(w, 400, introspectResponse)
	} else {
		log.Printf("Token is valid!")
		writeJSONResponse(w, 200, introspectResponse)
	}
	elapsed_s := float64(time.Since(start).Microseconds()) / 1000
	log.Printf("Elapsed register serve time - %.3f\n", elapsed_s)
	r.metricsClient.IncrementCounter(NUM_INTROSPECT, SERVICENAME)
	r.metricsClient.SetGaugeVal(ELAPSED_INTROSPECT_MS, SERVICENAME, elapsed_s)
}
