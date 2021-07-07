package metrics

import (
	"time"
)

// NopMetricsCollector implements the Collector interface with a set of
// Nop methods
type NopMetricsCollector struct {
}

func newNopMetricsCollector() Collector {
	return &NopMetricsCollector{}
}

func InitNopMetricsCollector() error {
	SetGlobMetricsCollector(newNopMetricsCollector())
	globalMetricsCollector.InitMetrics()
	return nil
}

// InitMetrics satisfies the Collector interface
func (NopMetricsCollector) InitMetrics() {}

// MeasureDurationMsSince satisfies the Collector interface
func (NopMetricsCollector) MeasureDurationMsSince(MetricName, time.Time, map[string]string) error {
	return nil
}

// MeasureDurationMs satisfies the Collector interface
func (NopMetricsCollector) MeasureDurationMs(MetricName, map[string]func() string) func() {
	return func() {}
}