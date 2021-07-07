package metrics

import (
	"bou.ke/monkey"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

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
	p.MeasureDurationMsSince("TEST_METRIC", histogramMap, starttime, labels)
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
	p.MeasureDurationMs("TEST_METRIC", histogramMap, map[string]func() string{})()
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