package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type API struct {
	Analyzer *Analyzer
	Redis    *redis.Client
}

func (a *API) Post(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ReqTotal.Inc()

	var s Sample
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.Analyzer.Push(s)

	data, _ := json.Marshal(s)
	a.Redis.Set(ctx, "last_metric", data, time.Minute)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("ok"))

	Latency.Observe(float64(time.Since(start).Milliseconds()))
}

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	res := a.Analyzer.Stats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
