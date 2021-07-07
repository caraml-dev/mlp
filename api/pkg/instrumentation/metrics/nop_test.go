package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"testing"
	"time"
)

// Nop methods return nothing and have no side effects.
// Simply exercise them to check that there are no panics.
func TestNopMethods(_ *testing.T) {
	testMetric := MetricName("TEST_METRIC")
	histogramMap := make(map[MetricName]*prometheus.HistogramVec)
	c := &NopMetricsCollector{}
	c.InitMetrics(histogramMap)
	c.MeasureDurationMs(testMetric, histogramMap, map[string]func() string{})
	c.MeasureDurationMsSince(testMetric, histogramMap, time.Now(), map[string]string{})
}