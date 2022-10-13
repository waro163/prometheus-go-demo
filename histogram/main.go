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

	// Create histogram, and only for the normal
	// distribution. The buckets are targeted to the parameters of the
	// normal distribution, with 20 buckets centered on the mean, each
	// half-sigma wide.
	hpptDurationsHistogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_request_durations_seconds",
		Help: "http request latency distributions.",
		// Buckets: prometheus.LinearBuckets(0.001, 1, 20),
		// Buckets: prometheus.ExponentialBuckets(0.001, 2, 20),
		Buckets: prometheus.ExponentialBucketsRange(1, 10, 10),
	})

	//
	hpptDurationsHistogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_durations_seconds_by_path",
		Help:    "http request latency distributions.",
		Buckets: []float64{0.001, 0.005, 0.01, 0.03, 0.05, 0.08, 0.1, 0.5, 0.7, 0.8, 1},
	},
		[]string{"service", "path"},
	)

	// register
	prometheus.MustRegister(hpptDurationsHistogram)
	prometheus.MustRegister(hpptDurationsHistogramVec)

	go func() {
		for {
			v := rand.Intn(10)
			hpptDurationsHistogram.Observe(float64(v))
			// inc totalRequest
			v1 := rand.Float64()
			if v < 3 {
				hpptDurationsHistogramVec.WithLabelValues("svc1", "/ping").Observe(v1)
			} else if v > 6 {
				hpptDurationsHistogramVec.With(prometheus.Labels{"service": "svc3", "path": "/health"}).Observe(v1)
			} else {
				hpptDurationsHistogramVec.WithLabelValues("svc2", "/probe").Observe(v1)
			}
			time.Sleep(time.Microsecond * 300)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
