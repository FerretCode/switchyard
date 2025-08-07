package prometheus

import (
	"github.com/ferretcode/switchyard/incident/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	loglineCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "loglines_ingested_total",
			Help: "Total number of ingested loglines",
		},
	)

	errorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "logline_errors_total",
			Help: "Total number of loglines ingested with ERROR severity",
		},
	)
)

func Init() types.PrometheusCounters {
	prometheus.MustRegister(loglineCounter)
	prometheus.MustRegister(errorCounter)

	return types.PrometheusCounters{
		LoglineCounter: loglineCounter,
		ErrorCounter:   errorCounter,
	}
}
