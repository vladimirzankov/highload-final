package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Routes(a *API) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/metrics", a.Post).Methods(http.MethodPost)
	r.HandleFunc("/analyze", a.Get).Methods(http.MethodGet)
	r.Handle("/metrics/prometheus", promhttp.Handler())

	return r
}
