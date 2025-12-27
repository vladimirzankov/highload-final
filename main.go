package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	cache := NewCache()

	analyzer := NewAnalyzer(50)
	analyzer.Start()

	prometheus.MustRegister(ReqTotal)
	prometheus.MustRegister(AnomalyTotal)
	prometheus.MustRegister(Latency)

	api := &API{
		Analyzer: analyzer,
		Redis:    cache,
	}

	r := Routes(api)

	log.Println("Service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
