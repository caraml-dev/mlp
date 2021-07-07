package metrics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	tu "github.com/gojek/mlp/api/internal/testutils"
)

func TestGetStatusString(t *testing.T) {
	assert.Equal(t, "success", GetStatusString(true))
	assert.Equal(t, "failure", GetStatusString(false))
}

func TestInitMetricsCollectorNop(t *testing.T) {
	err := InitMetricsCollector(false)
	// Validate
	assert.NoError(t, err)
	if _, ok := Glob().(*NopMetricsCollector); !ok {
		err := fmt.Errorf("Nop metrics collector was not initialised")
		tu.FailOnError(t, err)
	}
}

func TestInitMetricsCollectorPrometheus(t *testing.T) {
	err := InitMetricsCollector(true)
	// Validate
	assert.NoError(t, err)
	if _, ok := Glob().(*PrometheusClient); !ok {
		err := fmt.Errorf("Prometheus metrics collector was not initialised")
		tu.FailOnError(t, err)
	}
}
