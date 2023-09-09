package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	confidenceMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{ // INFO: ...Vec is a good idea, because that's effectively what labels do: They increase the number of time-series
		Name: "ml_simulation_confidence_percent",
		Help: "Confidence score for a simulated ML-Model inference",
	}, []string{"model"})
)

func simulateBinaryClassification(identifier string) {
	diminishingFactor := 1.0
	shouldRetrain := startRetrainWatchdog(&diminishingFactor)
	observer := confidenceMetric.WithLabelValues(identifier)

	for i := 0; ; i++ {
		// result will always be between 0.8 and 1.0
		result := (rand.Float64()*0.2 + 0.8) * diminishingFactor // In practice Prometheus always displays a value of approximately 0.9, given that the trend doesn't deteriorate

		observer.Observe(result)

		shouldRetrain <- result

		time.Sleep(time.Duration(rand.Float64()*900) * time.Millisecond)
		diminishingFactor -= 0.0001
	}
}

func startRetrainWatchdog(diminishingFactor *float64) chan float64 {
	shouldRetrain := make(chan float64)
	counter := 0

	go func() {
		for confidence := range shouldRetrain {
			if confidence < 0.7 {
				counter++
				if counter >= 50 {
					log.Println("Retraining!")
					*diminishingFactor = 1.0
					counter = 0
				}
			}
		}
	}()

	return shouldRetrain
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	for i := 1; i <= 5; i++ {
		log.Printf("Starting model %d\n", i)
		go simulateBinaryClassification(fmt.Sprint(i))
	}

	log.Fatal(http.ListenAndServe(":8080", nil)) // TODO: Port could be provided through helm / env-var
}
