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
	diminishingFactor = 1.0
	confidenceChannel = make(chan float64)
)

func simulateBinaryClassification() {
	confidenceMetric := promauto.NewSummaryVec(prometheus.SummaryOpts{ // INFO: ...Vec is a good idea, because that's effectively what labels do: They increase the number of time-series
		Name: "ml_simulation_confidence_percent",
		Help: "Confidence score for a simulated ML-Model inference",
	}, []string{"foo", "bar"}).
		// INFO: Normally you would set the labels not directly after the
		//  definition of a metric, but rather everywhere, where the labels
		//  stay the same. Defining a `prometheus.Observer` is more performant
		//  though. So it's generally a good advice to only observe on an
		//  Observer, rather than on setting the labels everytime.
		WithLabelValues("baz", "fizz")

	for {
		// result will always be between 0.8 and 1.0
		result := (rand.Float64()*0.2 + 0.8) * diminishingFactor // In practice Prometheus always displays a value of approximately 0.9, given that the trend doesn't deteriorate
		confidenceMetric.Observe(result)
		confidenceChannel <- result

		log.Println(result)

		time.Sleep(time.Duration(rand.Float64()*900) * time.Millisecond)
		diminishingFactor -= 0.0001
	}
}

func retrainListener() {
	counter := 0
	for confidence := range confidenceChannel {
		if confidence < 0.7 {
			counter++
			if counter >= 50 {
				log.Println("Retraining!")
				diminishingFactor = 1.0
				counter = 0
			}
		}
	}
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go simulateBinaryClassification()
	go retrainListener()
	http.ListenAndServe(":8080", nil)
}
