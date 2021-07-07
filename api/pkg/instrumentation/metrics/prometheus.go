package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/errgo.v2/fmt/errors"
)

// PrometheusHistogramVec is an interface that captures the methods from the the Prometheus
// HistogramVec type that are used in the app. This is added for unit testing.
type PrometheusHistogramVec interface {
	GetMetricWith(prometheus.Labels) (prometheus.Observer, error)
}

// getHistogramVec is a getter for the prometheus.HistogramVec defined for the input key.
// It returns a value satisfying the PrometheusHistogramVec interface
func getHistogramVec(key MetricName, histogramMap map[MetricName]*prometheus.HistogramVec) (PrometheusHistogramVec, error) {
	histVec, ok := histogramMap[key]
	if !ok {
		return nil, errors.Newf("Could not find the metric for %s", key)
	}
	return histVec, nil
}

// PrometheusClient satisfies the Collector interface
type PrometheusClient struct {
	histogramMap map[MetricName]*prometheus.HistogramVec
}

func InitPrometheusMetricsCollector(histogramMap map[MetricName]*prometheus.HistogramVec) error {
	SetGlobMetricsCollector(&PrometheusClient{histogramMap: histogramMap})
	globalMetricsCollector.InitMetrics()
	return nil
}

// InitMetrics initializes the collectors for all metrics defined for the app
// and registers them with the DefaultRegisterer.
func (p PrometheusClient) InitMetrics() {
	// Register histograms
	for _, obs := range p.histogramMap {
		prometheus.MustRegister(obs)
	}
}

// MeasureDurationMsSince takes in the Metric name, the start time and a map of labels and values
// to be associated to the metric. If errors occur in accessing the metric or associating the
// labels, they will simply be logged.
func (p PrometheusClient) MeasureDurationMsSince(
	key MetricName,
	starttime time.Time,
	labels map[string]string,
) error {
	// Get the histogram vec defined for the input key
	histVec, err := getHistogramVec(key, p.histogramMap)
	if err != nil {
		return err
	}
	// Create a histogram with the labels
	s, err := histVec.GetMetricWith(labels)
	if err != nil {
		return err
	}
	// Record the value in milliseconds
	s.Observe(float64(time.Since(starttime) / time.Millisecond))

	return nil
}

// MeasureDurationMs takes in the Metric name and a map of labels and functions to obtain
// the label values - this allows for MeasureDurationMs to be deferred and do a delayed
// evaluation of the labels. It returns a function which, when executed, will log the
// duration in ms since MeasureDurationMs was called. If errors occur in accessing the metric or
// associating the labels, they will simply be logged.
func (p PrometheusClient) MeasureDurationMs(
	key MetricName,
	labelValueGetters map[string]func() string,
) func() {
	// Capture start time
	starttime := time.Now()
	// Return function to measure and log the duration since start time
	return func() {
		// Evaluate the labels
		labels := map[string]string{}
		for key, f := range labelValueGetters {
			labels[key] = f()
		}
		// Log measurement
		p.MeasureDurationMsSince(key, starttime, labels)
	}
}
