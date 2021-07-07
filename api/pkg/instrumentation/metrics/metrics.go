package metrics

import (
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
	InitMetrics()
	MeasureDurationMsSince(key MetricName, starttime time.Time, labels map[string]string) error
	// MeasureDurationMs is a deferrable version of MeasureDurationMsSince which evaluates labels
	// at the time of logging
	MeasureDurationMs(key MetricName, labels map[string]func() string) func()
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

// GetStatusString returns a classification string (success / failure) based on the
// input boolean
func GetStatusString(status bool) string {
	if status {
		return statusLabels.Success
	}
	return statusLabels.Failure
}
