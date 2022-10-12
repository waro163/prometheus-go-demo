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

	// gauge
	totalUsage := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_total",
			Help: "calculate total memory usage",
		},
	)
	// gaugevec
	podUsage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_usage_by_pod",
			Help: "calculate each pod memory usage",
		},
		[]string{"host", "pod"},
	)

	// register
	prometheus.MustRegister(totalUsage)
	prometheus.MustRegister(podUsage)

	podUsage.WithLabelValues("vm1", "pod01").Set(1280)
	go func() {
		for {
			v := rand.Intn(10)
			if v%5 == 0 {
				podUsage.With(prometheus.Labels{"host": "vm5", "pod": "pod03"}).Set(float64(v))
			} else {
				podUsage.With(prometheus.Labels{"host": "vm5", "pod": "pod01"}).Add(float64(v))
			}
			if v%3 == 0 {
				podUsage.With(prometheus.Labels{"host": "vm3", "pod": "pod01"}).Set(float64(v))
			} else {
				podUsage.With(prometheus.Labels{"host": "vm3", "pod": "pod02"}).Inc()
			}
			// inc totalRequest
			totalUsage.Add(float64(v))

			time.Sleep(time.Microsecond * 300)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
