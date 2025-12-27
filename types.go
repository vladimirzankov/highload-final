package main

type Sample struct {
	Timestamp int64   `json:"timestamp"`
	CPU       float64 `json:"cpu"`
	RPS       float64 `json:"rps"`
}

type Report struct {
	RollingAvg float64 `json:"avg"`
	ZScore     float64 `json:"zscore"`
	Anomaly    bool    `json:"anomaly"`
}
