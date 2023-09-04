package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	confidence = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "ml_simulation_confidence_percent",
		Help: "Confidence score for a simulated ML-Model inference",
	})
)

func simulateBinaryClassification() {
	for {
		// result will always be between 0.8 and 1.0
		result := rand.Float64()*0.2 + 0.8 // In practice Prometheus always displays a value of approximately 0.9, given that the trend doesn't deteriorate
		confidence.Observe(result)
		log.Println(result)
		time.Sleep(time.Duration(rand.Float64()*900) * time.Millisecond)
	}
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go simulateBinaryClassification()
	http.ListenAndServe(":8080", nil)
}
