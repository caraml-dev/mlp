package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// NopMetricsCollector implements the Collector interface with a set of
// Nop methods
type NopMetricsCollector struct {
}

func newNopMetricsCollector() Collector {
	return &NopMetricsCollector{}
}

// InitMetrics satisfies the Collector interface
func (NopMetricsCollector) InitMetrics(map[MetricName]*prometheus.HistogramVec) {}

// MeasureDurationMsSince satisfies the Collector interface
func (NopMetricsCollector) MeasureDurationMsSince(MetricName, map[MetricName]*prometheus.HistogramVec, time.Time, map[string]string) {}

// MeasureDurationMs satisfies the Collector interface
func (NopMetricsCollector) MeasureDurationMs(MetricName, map[MetricName]*prometheus.HistogramVec, map[string]func() string) func() {
	return func() {}
}