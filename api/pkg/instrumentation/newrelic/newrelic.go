package newrelic

import (
	"fmt"
	"net/http"
	"time"

	newrelic "github.com/newrelic/go-agent"
)

var (
	newRelicApp newrelic.Application = &NoopApp{}
)

// Config stores NewRelic configuration.
type Config struct {
	Enabled bool                   `envconfig:"NEWRELIC_ENABLED" default:"false"`
	AppName string                 `envconfig:"NEWRELIC_APP_NAME" default:"app-name"`
	License string                 `envconfig:"NEWRELIC_LICENSE" default:""`
	Labels  map[string]interface{} `envconfig:"NEWRELIC_LABELS" default:""`
	// https://docs.newrelic.com/docs/agents/go-agent/configuration/go-agent-configuration#error-ignore-status
	IgnoreStatusCodes []int `envconfig:"NEWRELIC_IGNORE_STATUS_CODES" default:"400,401,403,404,405,412"`
}

// InitNewRelic initializes NewRelic Application.
func InitNewRelic(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}

	if cfg.License == "" {
		return nil
	}

	config := newrelic.NewConfig(cfg.AppName, cfg.License)
	for k, v := range cfg.Labels {
		config.Labels[k] = fmt.Sprint(v)
	}
	config.ErrorCollector.IgnoreStatusCodes = cfg.IgnoreStatusCodes

	app, err := newrelic.NewApplication(config)
	if err != nil {
		return err
	}

	newRelicApp = app
	return nil
}

// StartTransaction implements newrelic.Application interface.
func StartTransaction(name string, w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	return newRelicApp.StartTransaction(name, w, r)
}

// RecordCustomEvent implements newrelic.Application interface.
func RecordCustomEvent(eventType string, params map[string]interface{}) error {
	return newRelicApp.RecordCustomEvent(eventType, params)
}

// RecordCustomMetric implements newrelic.Application interface.
func RecordCustomMetric(name string, value float64) error {
	return newRelicApp.RecordCustomMetric(name, value)
}

// Shutdown flushes data to New Relic's servers and stops all
// agent-related goroutines managing this application.
func Shutdown(timeout time.Duration) {
	newRelicApp.Shutdown(timeout)
}

// WrapHandleFunc wraps newrelic.WrapHandleFunc.
func WrapHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	return newrelic.WrapHandleFunc(newRelicApp, pattern, handler)
}

// WrapHandle wraps newrelic.WrapHandle.
func WrapHandle(pattern string, handler http.Handler) (string, http.Handler) {
	return newrelic.WrapHandle(newRelicApp, pattern, handler)
}
