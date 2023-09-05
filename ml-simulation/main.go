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
	labelValues       = []string{"fizz", "buzz"}
)

func simulateBinaryClassification() {
	confidenceMetric := promauto.NewSummaryVec(prometheus.SummaryOpts{ // INFO: ...Vec is a good idea, because that's effectively what labels do: They increase the number of time-series
		Name: "ml_simulation_confidence_percent",
		Help: "Confidence score for a simulated ML-Model inference",
	}, []string{"foo", "bar"})

	for i := 0; ; i++ {
		// result will always be between 0.8 and 1.0
		result := (rand.Float64()*0.2 + 0.8) * diminishingFactor // In practice Prometheus always displays a value of approximately 0.9, given that the trend doesn't deteriorate

		// INFO: This just constructs different time-series depending on some random condition
		if j := i % 2; j > 0 {
			confidenceMetric.WithLabelValues(labelValues[0], labelValues[1]).Observe(result)
		} else {
			confidenceMetric.WithLabelValues(labelValues[1], labelValues[0]).Observe(result)
		}

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
