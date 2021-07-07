package metrics

import (
	"fmt"
	"github.com/gojek/mlp/api/log"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/errgo.v2/fmt/errors"
	"time"
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
		fmt.Println("================ failed ====================")
		return nil, errors.Newf("Could not find the metric for %s", key)
	}
	return histVec, nil
}

// PrometheusClient satisfies the Collector interface
type PrometheusClient struct {
}

// InitMetrics initializes the collectors for all metrics defined for the app
// and registers them with the DefaultRegisterer.
func (PrometheusClient) InitMetrics(histogramMap map[MetricName]*prometheus.HistogramVec) {
	// Register histograms
	for _, obs := range histogramMap {
		prometheus.MustRegister(obs)
	}
}

// MeasureDurationMsSince takes in the Metric name, the start time and a map of labels and values
// to be associated to the metric. If errors occur in accessing the metric or associating the
// labels, they will simply be logged.
func (PrometheusClient) MeasureDurationMsSince(
	key MetricName,
	histogramMap map[MetricName]*prometheus.HistogramVec,
	starttime time.Time,
	labels map[string]string,
) {
	// Get the histogram vec defined for the input key
	histVec, err := getHistogramVec(key, histogramMap)
	if err != nil {
		log.GlobalLogger.Errorf(err.Error())
		return
	}
	// Create a histogram with the labels
	s, err := histVec.GetMetricWith(labels)
	if err != nil {
		log.GlobalLogger.Errorf("Error occurred when creating histogram for %s: %v", key, err)
		return
	}
	// Record the value in milliseconds
	s.Observe(float64(time.Since(starttime) / time.Millisecond))
}

// MeasureDurationMs takes in the Metric name and a map of labels and functions to obtain
// the label values - this allows for MeasureDurationMs to be deferred and do a delayed
// evaluation of the labels. It returns a function which, when executed, will log the
// duration in ms since MeasureDurationMs was called. If errors occur in accessing the metric or
// associating the labels, they will simply be logged.
func (p PrometheusClient) MeasureDurationMs(
	key MetricName,
	histogramMap map[MetricName]*prometheus.HistogramVec,
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
		p.MeasureDurationMsSince(key, histogramMap, starttime, labels)
	}
}
