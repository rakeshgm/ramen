package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	metricNamespace = "ramen"
)

var (
	metricLabels = []string{
		"resource_type",     // Name of the resource [drpc|vrg]
		"name",              // Name of the resource [drpc-name|vrg-name]
		"namespace",         // drpc namespace name
		"schedule_interval", // value from drpolicy
	}

	lastSyncTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "last_sync_timestamp_seconds",
			Namespace: metricNamespace,
			Help:      "Placeholder text",
		},
		metricLabels,
	)
)

func NewLastSyncMetrics(lables prometheus.Labels) prometheus.Gauge {
	return lastSyncTime.With(lables)
}

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(lastSyncTime)
}
