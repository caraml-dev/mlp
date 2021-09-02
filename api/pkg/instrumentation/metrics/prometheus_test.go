package metrics

import (
	"testing"
	"time"

	io_prometheus_client "github.com/prometheus/client_model/go"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type noOpCollector struct{}

func (c *noOpCollector) Describe(chan<- *prometheus.Desc) {}

func (c *noOpCollector) Collect(chan<- prometheus.Metric) {}

type mockCounter struct {
	mock.Mock
	noOpCollector
}

func (c *mockCounter) Desc() *prometheus.Desc {
	c.Called()
	return nil
}

func (c *mockCounter) Write(metric *io_prometheus_client.Metric) error {
	c.Called(metric)
	return nil
}

func (c *mockCounter) Inc() {
	c.Called()
}

func (c *mockCounter) Add(f float64) {
	c.Called(f)
}

type mockCounterVec struct {
	noOpCollector
	counter *mockCounter
}

func (m mockCounterVec) GetMetricWith(labels prometheus.Labels) (prometheus.Counter, error) {
	return m.counter, nil
}

// mockGauge mocks a prometheus Gauge
type mockGauge struct {
	mock.Mock
	noOpCollector
	value float64
}

// Implementing Prometheus Gauge interface
func (g *mockGauge) Desc() *prometheus.Desc {
	return nil
}

func (g *mockGauge) Write(*io_prometheus_client.Metric) error {
	return nil
}

func (g *mockGauge) SetToCurrentTime() {
	g.Called()
}

func (g *mockGauge) Set(value float64) {
	g.Called(value)
}

func (g *mockGauge) Inc() {
	g.Add(1)
}

func (g *mockGauge) Dec() {
	g.Sub(1)
}

func (g *mockGauge) Add(value float64) {
	g.Called(value)
}

func (g *mockGauge) Sub(value float64) {
	g.Called(value)
}

// mockGaugeVec mocks a prometheus GaugeVec
type mockGaugeVec struct {
	noOpCollector
	gauge *mockGauge
}

func (g *mockGaugeVec) GetMetricWith(labels prometheus.Labels) (prometheus.Gauge, error) {
	return g.gauge, nil
}

// createMockGaugeVec creates a mock gauge and a mock gauge vec
func createMockGaugeVec(testValue float64) *mockGaugeVec {
	// Create mock gauge and gauge vec
	gauge := &mockGauge{
		value: 0,
	}
	gauge.On("Set", mock.Anything).Run(func(args mock.Arguments) {
		gauge.value = testValue
	}).Return(nil)
	gaugeVec := &mockGaugeVec{
		gauge: gauge,
	}
	return gaugeVec
}

var gaugeMap = make(map[MetricName]PrometheusGaugeVec)

func TestGetGaugeVec(t *testing.T) {
	_, err := getGaugeVec("TEST_METRIC", gaugeMap)
	assert.Error(t, err)
}

func TestMeasureGauge(t *testing.T) {
	metricName := MetricName("TEST_METRIC")
	value := float64(5)
	labels := map[string]string{}
	// Create mock gauge vec
	gaugeVec := createMockGaugeVec(value)
	p := &PrometheusClient{
		gaugeMap: map[MetricName]PrometheusGaugeVec{
			metricName: gaugeVec,
		},
	}

	err := p.RecordGauge(metricName, value, labels)
	assert.NoError(t, err)
	gaugeVec.gauge.AssertCalled(t, "Set", mock.AnythingOfType("float64"))
	assert.Equal(t, value, gaugeVec.gauge.value)
}

// mockHistogramVec mocks a prometheus HistogramVec
type mockHistogramVec struct {
	noOpCollector
	histogram *mockHistogram
}

func (h *mockHistogramVec) GetMetricWith(labels prometheus.Labels) (prometheus.Observer, error) {
	return h.histogram, nil
}

// mockHistogram mocks a prometheus Histogram
type mockHistogram struct {
	mock.Mock
	duration float64
}

// Implementing Prometheus Histogram interface
func (h *mockHistogram) Observe(duration float64) {
	h.Called(duration)
}

var histogramMap = make(map[MetricName]PrometheusHistogramVec)

func TestGetHistogramVec(t *testing.T) {
	_, err := getHistogramVec("TEST_METRIC", histogramMap)
	assert.Error(t, err)
}

func TestMeasureDurationMsSince(t *testing.T) {
	metricName := MetricName("TEST_METRIC")
	starttime := time.Now()
	labels := map[string]string{}
	testDuration := 100.0
	// Create mock histogram vec
	histVec := createMockHistVec(testDuration)
	p := &PrometheusClient{
		histogramMap: map[MetricName]PrometheusHistogramVec{
			metricName: histVec,
		},
	}
	err := p.MeasureDurationMsSince(metricName, starttime, labels)
	assert.NoError(t, err)
	histVec.histogram.AssertCalled(t, "Observe", mock.AnythingOfType("float64"))
	assert.Equal(t, testDuration, histVec.histogram.duration)
}

func TestMeasureDurationMs(t *testing.T) {
	metricName := MetricName("TEST_METRIC")
	testDuration := 200.0
	// Create mock histogram vec
	histVec := createMockHistVec(testDuration)
	p := &PrometheusClient{
		histogramMap: map[MetricName]PrometheusHistogramVec{
			metricName: histVec,
		},
	}
	p.MeasureDurationMs(metricName, map[string]func() string{})()
	// Validate
	histVec.histogram.AssertCalled(t, "Observe", mock.AnythingOfType("float64"))
	assert.Equal(t, testDuration, histVec.histogram.duration)
}

// createMockHistVec creates a mock histogram and a mock histogram vec
func createMockHistVec(testDuration float64) *mockHistogramVec {
	// Create mock histogram and histogram vec
	hist := &mockHistogram{
		duration: 0,
	}
	hist.On("Observe", mock.Anything).Run(func(args mock.Arguments) {
		hist.duration = testDuration
	}).Return(nil)
	return &mockHistogramVec{
		histogram: hist,
	}
}

func TestCounterInc(t *testing.T) {
	metricName := MetricName("TEST_METRIC")
	counter := &mockCounter{}
	counter.On("Inc").Return()
	counterVec := &mockCounterVec{counter: counter}
	p := &PrometheusClient{
		counterMap: map[MetricName]PrometheusCounterVec{
			metricName: counterVec,
		},
	}
	p.Inc(metricName, map[string]string{})
	counter.AssertCalled(t, "Inc")
}
