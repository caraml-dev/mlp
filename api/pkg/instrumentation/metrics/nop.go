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

	return nil
}

// MeasureDurationMsSince satisfies the Collector interface
func (NopMetricsCollector) MeasureDurationMsSince(MetricName, time.Time, map[string]string) error {
	return nil
}

// MeasureDurationMs satisfies the Collector interface
func (NopMetricsCollector) MeasureDurationMs(MetricName, map[string]func() string) func() {
	return func() {}
}

// RecordGauge satisfies the Collector interface
func (NopMetricsCollector) RecordGauge(MetricName, float64, map[string]string) error {
	return nil
}

// Inc satisfies the Collector interface
func (c NopMetricsCollector) Inc(_ MetricName, _ map[string]string) error {
	return nil
}
