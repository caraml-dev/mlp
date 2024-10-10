package newrelic

import (
	"net/http"
	"time"

	newrelic "github.com/newrelic/go-agent"
)

// A NoopApp is a trivial, minimum overhead implementation of newrelic.Application
// for which all operations are no-ops.
type NoopApp struct{}

// StartTransaction implements newrelic.Application interface.
func (na NoopApp) StartTransaction(_ string, w http.ResponseWriter, _ *http.Request) newrelic.Transaction {
	return &NoopTx{
		w: w,
	}
}

// RecordCustomEvent implements newrelic.Application interface.
func (na NoopApp) RecordCustomEvent(_ string, _ map[string]interface{}) error {
	return nil
}

// RecordCustomMetric implements newrelic.Application interface.
func (na NoopApp) RecordCustomMetric(_ string, _ float64) error { return nil }

// WaitForConnection implements newrelic.Application interface.
func (na NoopApp) WaitForConnection(_ time.Duration) error { return nil }

// Shutdown implements newrelic.Application interface.
func (na NoopApp) Shutdown(_ time.Duration) {
	// Do nothing
}

// A NoopTx is a trivial, minimum overhead implementation of newrelic.Transaction
// for which all operations are no-ops.
type NoopTx struct {
	w http.ResponseWriter
}

func (nt *NoopTx) IsSampled() bool {
	return false
}

// End implements newrelic.Transaction interface.
func (nt *NoopTx) End() error {
	return nil
}

// Ignore implements newrelic.Transaction interface.
func (nt *NoopTx) Ignore() error {
	return nil
}

// SetName implements newrelic.Transaction interface.
func (nt *NoopTx) SetName(_ string) error {
	return nil
}

// NoticeError implements newrelic.Transaction interface.
func (nt *NoopTx) NoticeError(_ error) error {
	return nil
}

// AddAttribute implements newrelic.Transaction interface.
func (nt *NoopTx) AddAttribute(_ string, _ interface{}) error {
	return nil
}

// SetWebRequest implements newrelic.Transaction interface.
func (nt *NoopTx) SetWebRequest(newrelic.WebRequest) error {
	return nil
}

// SetWebResponse implements newrelic.Transaction interface.
func (nt *NoopTx) SetWebResponse(http.ResponseWriter) newrelic.Transaction {
	return nil
}

// StartSegmentNow implements newrelic.Transaction interface.
func (nt *NoopTx) StartSegmentNow() newrelic.SegmentStartTime {
	return newrelic.SegmentStartTime{}
}

// CreateDistributedTracePayload implements newrelic.Transaction interface.
func (nt *NoopTx) CreateDistributedTracePayload() newrelic.DistributedTracePayload {
	return nil
}

// AcceptDistributedTracePayload implements newrelic.Transaction interface.
func (nt *NoopTx) AcceptDistributedTracePayload(_ newrelic.TransportType, _ interface{}) error {
	return nil
}

// Application implements newrelic.Transaction interface.
func (nt *NoopTx) Application() newrelic.Application {
	return nil
}

// BrowserTimingHeader implements newrelic.Transaction interface.
func (nt *NoopTx) BrowserTimingHeader() (*newrelic.BrowserTimingHeader, error) {
	return nil, nil
}

// NewGoroutine implements newrelic.Transaction interface.
func (nt *NoopTx) NewGoroutine() newrelic.Transaction {
	return nil
}

// Header implements http.ResponseWriter interface.
func (nt *NoopTx) Header() http.Header {
	return nt.w.Header()
}

// Write implements http.ResponseWriter interface.
func (nt *NoopTx) Write(b []byte) (int, error) {
	return nt.w.Write(b)
}

// WriteHeader implements http.ResponseWriter interface.
func (nt *NoopTx) WriteHeader(code int) {
	nt.w.WriteHeader(code)
}

// GetTraceMetadata implements newrelic.Transaction interface.
func (nt *NoopTx) GetTraceMetadata() newrelic.TraceMetadata {
	return newrelic.TraceMetadata{}
}

// GetLinkingMetadata implements newrelic.Transaction interface.
func (nt *NoopTx) GetLinkingMetadata() newrelic.LinkingMetadata {
	return newrelic.LinkingMetadata{}
}
