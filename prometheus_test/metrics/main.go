package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	myCounter          = prometheus.NewCounter(prometheus.CounterOpts{Name: "myCounter", Help: "Test Counter"})
	myGauge            = prometheus.NewGauge(prometheus.GaugeOpts{Name: "myGauge", Help: "Test Gauge"})
	mySummary          = prometheus.NewSummary(prometheus.SummaryOpts{Name: "mySummary", Help: "Test Summary", Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}})
	myHistogram        = prometheus.NewHistogram(prometheus.HistogramOpts{Name: "myHistogram", Help: "Test Histogram", Buckets: prometheus.LinearBuckets(95, 1, 10)})
	myCounterWithLabel = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "myCounterWithLabel", Help: "Test Counter", ConstLabels: prometheus.Labels{"name": "xxxx"}}, []string{"type"})
)

func init() {
	prometheus.MustRegister(
		myCounter,
		myCounterWithLabel,
		myGauge,
		mySummary,
		myHistogram,
	)
}

func main() {

	go func() {
		for {
			myCounter.Inc()
			myCounterWithLabel.WithLabelValues("type_1").Inc()
			myCounterWithLabel.WithLabelValues("type_2").Inc()
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			myGauge.Set(float64(100 + rand.Int31n(10) - 5))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			mySummary.Observe(float64(100 + rand.Int31n(10) - 5))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			myHistogram.Observe(float64(100 + rand.Int31n(10) - 5))
			time.Sleep(1 * time.Second)
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(":8889", nil))
}
