package metrics

import (
	"github.com/gojek/mlp/api/log"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// MetricName is a type used to define the names of the various metrics collected in the App.
type MetricName string

var statusLabels = struct {
	Success string
	Failure string
}{
	Success: "success",
	Failure: "failure",
}

// Collector defines the common interface for all metrics collection engines
type Collector interface {
	InitMetrics(histogramMap map[MetricName]*prometheus.HistogramVec)
	MeasureDurationMsSince(key MetricName, histogramMap map[MetricName]*prometheus.HistogramVec, starttime time.Time, labels map[string]string)
	// MeasureDurationMs is a deferrable version of MeasureDurationMsSince which evaluates labels
	// at the time of logging
	MeasureDurationMs(key MetricName, histogramMap map[MetricName]*prometheus.HistogramVec, labels map[string]func() string) func()
}

// globalMetricsCollector is initialised to a Nop metrics collector. Calling
// InitMetricsCollector can update this value.
var globalMetricsCollector = newNopMetricsCollector()

// Glob returns the global metrics collector
func Glob() Collector {
	return globalMetricsCollector
}

// SetGlobMetricsCollector is used to update the global metrics collector instance with the input
func SetGlobMetricsCollector(c Collector) {
	globalMetricsCollector = c
}

// InitMetricsCollector is used to select the appropriate metrics collector and
// set up the required values for instrumenting.
func InitMetricsCollector(enabled bool) error {
	if enabled {
		log.GlobalLogger.Info("Initializing Prometheus Metrics Collector")
		// Use the Prometheus Instrumentation Client
		SetGlobMetricsCollector(&PrometheusClient{})
	} else {
		// Use the Nop Metrics collector
		log.GlobalLogger.Info("Initializing Nop Metrics Collector")
		SetGlobMetricsCollector(newNopMetricsCollector())
	}
	// Initialize
	histogramMap := make(map[MetricName]*prometheus.HistogramVec)
	globalMetricsCollector.InitMetrics(histogramMap)
	return nil
}

// GetStatusString returns a classification string (success / failure) based on the
// input boolean
func GetStatusString(status bool) string {
	if status {
		return statusLabels.Success
	}
	return statusLabels.Failure
}