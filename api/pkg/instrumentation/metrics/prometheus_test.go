package metrics

import (
	io_prometheus_client "github.com/prometheus/client_model/go"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockGauge mocks a prometheus Gauge
type mockGauge struct {
	mock.Mock
	value float64
}

// Implementing Prometheus Gauge interface
func (g *mockGauge) Desc() *prometheus.Desc {
	return nil
}

func (g *mockGauge) Write(*io_prometheus_client.Metric) error {
	return nil
}

func (g *mockGauge) Describe(chan<- *prometheus.Desc) {}

func (g *mockGauge) Collect(chan<- prometheus.Metric) {}

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
	mock.Mock
	gauge *mockGauge
}

func (g *mockGaugeVec) GetMetricWith(labels prometheus.Labels) (prometheus.Gauge, error) {
	g.Called(labels)
	// Return mockGauge
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
	gaugeVec.On("GetMetricWith", mock.Anything).Return(gauge, nil)
	return gaugeVec
}

var gaugeMap = make(map[MetricName]*prometheus.GaugeVec)

func TestGetGaugeVec(t *testing.T) {
	_, err := getGaugeVec("TEST_METRIC", gaugeMap)
	assert.Error(t, err)
}

func TestMeasureGauge(t *testing.T) {
	p := &PrometheusClient{}
	value := float64(5)
	labels := map[string]string{}
	// Create mock gauge vec
	gaugeVec := createMockGaugeVec(value)
	// Patch getGaugeVec for the test and run
	monkey.Patch(getGaugeVec,
		func(key MetricName, gaugeMap map[MetricName]*prometheus.GaugeVec) (PrometheusGaugeVec, error) {
			return gaugeVec, nil
		})
	p.RecordGauge("TEST_METRIC", value, labels)
	monkey.Unpatch(getGaugeVec)
	// Validate
	gaugeVec.AssertCalled(t, "GetMetricWith", mock.Anything)
	gaugeVec.gauge.AssertCalled(t, "Set", mock.AnythingOfType("float64"))
	assert.Equal(t, value, gaugeVec.gauge.value)
}

// mockHistogramVec mocks a prometheus HistogramVec
type mockHistogramVec struct {
	mock.Mock
	histogram *mockHistogram
}

func (h *mockHistogramVec) GetMetricWith(labels prometheus.Labels) (prometheus.Observer, error) {
	h.Called(labels)
	// Return mockHistogram
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

var histogramMap = make(map[MetricName]*prometheus.HistogramVec)

func TestGetHistogramVec(t *testing.T) {
	_, err := getHistogramVec("TEST_METRIC", histogramMap)
	assert.Error(t, err)
}

func TestMeasureDurationMsSince(t *testing.T) {
	p := &PrometheusClient{}
	starttime := time.Now()
	labels := map[string]string{}
	testDuration := 100.0
	// Create mock histogram vec
	histVec := createMockHistVec(testDuration)
	// Patch getHistogramVec for the test and run
	monkey.Patch(getHistogramVec,
		func(key MetricName, histogramMap map[MetricName]*prometheus.HistogramVec) (PrometheusHistogramVec, error) {
			return histVec, nil
		})
	p.MeasureDurationMsSince("TEST_METRIC", starttime, labels)
	monkey.Unpatch(getHistogramVec)
	// Validate
	histVec.AssertCalled(t, "GetMetricWith", mock.Anything)
	histVec.histogram.AssertCalled(t, "Observe", mock.AnythingOfType("float64"))
	assert.Equal(t, testDuration, histVec.histogram.duration)
}

func TestMeasureDurationMs(t *testing.T) {
	p := &PrometheusClient{}
	testDuration := 200.0
	// Create mock histogram vec
	histVec := createMockHistVec(testDuration)
	// Patch getHistogramVec for the test and run
	monkey.Patch(getHistogramVec,
		func(key MetricName, histogramMap map[MetricName]*prometheus.HistogramVec) (PrometheusHistogramVec, error) {
			return histVec, nil
		})
	p.MeasureDurationMs("TEST_METRIC", map[string]func() string{})()
	monkey.Unpatch(getHistogramVec)
	// Validate
	histVec.AssertCalled(t, "GetMetricWith", mock.Anything)
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
	histVec := &mockHistogramVec{
		histogram: hist,
	}
	histVec.On("GetMetricWith", mock.Anything).Return(hist, nil)
	return histVec
}
