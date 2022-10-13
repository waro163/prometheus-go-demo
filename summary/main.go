// A minimal example of how to include Prometheus instrumentation.
package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	rand.Seed(time.Now().Unix())

	// Create
	hpptDurationsSummary := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "http_request_durations_seconds_summary",
		Help:       "http request latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})

	//
	hpptDurationsSummaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_request_durations_seconds_summary_by_service",
		Help:       "http request latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
		[]string{"service"},
	)

	// register
	prometheus.MustRegister(hpptDurationsSummary)
	prometheus.MustRegister(hpptDurationsSummaryVec)

	go func() {
		for {
			v := rand.Intn(10)
			hpptDurationsSummary.Observe(float64(v))
			// inc totalRequest
			v1 := rand.Float64()
			if v < 3 {
				hpptDurationsSummaryVec.WithLabelValues("uniform").Observe(v1)
			} else if v > 6 {
				hpptDurationsSummaryVec.With(prometheus.Labels{"service": "normal"}).Observe(v1)
			} else {
				hpptDurationsSummaryVec.WithLabelValues("exponential").Observe(v1)
			}
			time.Sleep(time.Microsecond * 300)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
