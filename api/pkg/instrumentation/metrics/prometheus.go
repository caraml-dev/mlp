package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/errgo.v2/fmt/errors"
)

// PrometheusGaugeVec is an interface that captures the methods from the the Prometheus
// GaugeVec type that are used in the app. This is added for unit testing.
type PrometheusGaugeVec interface {
	GetMetricWith(prometheus.Labels) (prometheus.Gauge, error)
}

// getGaugeVec is a getter for the prometheus.GaugeVec defined for the input key.
// It returns a value satisfying the PrometheusGaugeVec interface
func getGaugeVec(key MetricName, gaugeMap map[MetricName]*prometheus.GaugeVec) (PrometheusGaugeVec, error) {
	gaugeVec, ok := gaugeMap[key]
	if !ok {
		return nil, errors.Newf("Could not find the metric for %s", key)
	}
	return gaugeVec, nil
}

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
	gaugeMap map[MetricName]*prometheus.GaugeVec
	histogramMap map[MetricName]*prometheus.HistogramVec
}

// InitPrometheusMetricsCollector initializes the collectors for all metrics defined for the app
// and registers them with the DefaultRegisterer.
func InitPrometheusMetricsCollector(
	gaugeMap map[MetricName]*prometheus.GaugeVec,
	histogramMap map[MetricName]*prometheus.HistogramVec,
) error {
	SetGlobMetricsCollector(&PrometheusClient{
		gaugeMap: gaugeMap,
		histogramMap: histogramMap,
	})
	for _, obs := range gaugeMap {
		prometheus.MustRegister(obs)
	}
	for _, obs := range histogramMap {
		prometheus.MustRegister(obs)
	}

	return nil
}

// MeasureDurationMsSince takes in the Metric name, the start time and a map of labels and values
// to be associated to the metric. Error will be thrown if errors occur in accessing
// the metric or associating the labels.
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

// RecordGauge takes in the Metric name, the count and a map of labels and values
// to be associated to the metric. Error will be thrown if errors occur in accessing
// the metric or associating the labels.
func (p PrometheusClient) RecordGauge(
	key MetricName,
	value float64,
	labels map[string]string,
) error {
	// Get the gauge vec defined for the input key
	gaugeVec, err := getGaugeVec(key, p.gaugeMap)
	if err != nil {
		return err
	}
	// Create a gauge with the labels
	s, err := gaugeVec.GetMetricWith(labels)
	if err != nil {
		return err
	}

	// Record the value
	s.Set(value)

	return nil
}
