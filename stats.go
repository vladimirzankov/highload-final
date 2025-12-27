package main

import (
	"math"
	"sync"
)

type Analyzer struct {
	size   int
	buf    []float64
	report Report
	mu     sync.Mutex
	ch     chan Sample
}

func NewAnalyzer(size int) *Analyzer {
	return &Analyzer{
		size: size,
		buf:  make([]float64, 0, size),
		ch:   make(chan Sample, 10000),
	}
}

func (a *Analyzer) Start() {
	go func() {
		for s := range a.ch {
			a.update(s)
		}
	}()
}

func (a *Analyzer) Push(s Sample) {
	select {
	case a.ch <- s:
	default:
	}
}

func (a *Analyzer) Stats() Report {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.report
}

func (a *Analyzer) update(s Sample) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.buf = append(a.buf, s.RPS)
	if len(a.buf) > a.size {
		a.buf = a.buf[1:]
	}

	m := avg(a.buf)
	d := sd(a.buf, m)

	z := 0.0
	if d > 0 {
		z = (s.RPS - m) / d
	}

	anomaly := math.Abs(z) > 2.0
	if anomaly {
		AnomalyTotal.Inc()
	}

	a.report = Report{
		RollingAvg: m,
		ZScore:     z,
		Anomaly:    anomaly,
	}
}

func avg(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func sd(data []float64, mean float64) float64 {
	var sum float64
	for _, v := range data {
		diff := v - mean
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(data)))
}
