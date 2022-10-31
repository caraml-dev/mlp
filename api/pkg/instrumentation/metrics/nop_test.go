package metrics

import (
	"testing"
	"time"
)

// Nop methods return nothing and have no side effects.
// Simply exercise them to check that there are no panics.
func TestNopMethods(_ *testing.T) {
	testMetric := MetricName("TEST_METRIC")
	c := &NopMetricsCollector{}
	c.MeasureDurationMs(testMetric, map[string]func() string{})
	_ = c.MeasureDurationMsSince(testMetric, time.Now(), map[string]string{})
	_ = c.RecordGauge(testMetric, 0, map[string]string{})
}
