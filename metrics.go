package main

import "github.com/prometheus/client_golang/prometheus"

var (
	ReqTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
	)

	AnomalyTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "anomalies_total",
			Help: "Total detected anomalies",
		},
	)

	Latency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_latency_ms",
			Help:    "HTTP request latency in milliseconds",
			Buckets: prometheus.LinearBuckets(5, 5, 10),
		},
	)
)
