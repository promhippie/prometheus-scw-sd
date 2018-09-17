package action

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/promhippie/prometheus-scw-sd/pkg/version"
)

var (
	registry  = prometheus.NewRegistry()
	namespace = "prometheus_scw_sd"
)

var (
	requestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "request_duration_seconds",
			Help:      "Histogram of latencies for requests to the Scaleway API.",
			Buckets:   []float64{0.001, 0.01, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0},
		},
	)

	requestFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "request_failures_total",
			Help:      "Total number of failed requests to the Scaleway API.",
		},
	)
)

func init() {
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
	registry.MustRegister(prometheus.NewGoCollector())

	registry.MustRegister(version.Collector(namespace))

	registry.MustRegister(requestDuration)
	registry.MustRegister(requestFailures)
}

type promLogger struct {
	logger log.Logger
}

func (pl promLogger) Println(v ...interface{}) {
	level.Error(pl.logger).Log(
		"msg", fmt.Sprintln(v...),
	)
}
