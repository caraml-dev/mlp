package newrelic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	newrelic "github.com/newrelic/go-agent"
)

func TestNoopApp(t *testing.T) {
	na := NoopApp{}
	na.StartTransaction("test", httptest.NewRecorder(), &http.Request{})
	na.RecordCustomEvent("test", nil)
	na.RecordCustomMetric("test", 0)
	na.WaitForConnection(0)
	na.Shutdown(0)
}

func TestNoopTx(t *testing.T) {
	nt := NoopTx{
		w: httptest.NewRecorder(),
	}
	nt.End()
	nt.Ignore()
	nt.SetName("test")
	nt.NoticeError(nil)
	nt.AddAttribute("key", "val")
	nt.SetWebRequest(nil)
	nt.SetWebResponse(nil)
	nt.StartSegmentNow()
	nt.CreateDistributedTracePayload()
	nt.AcceptDistributedTracePayload(newrelic.TransportUnknown, nil)
	nt.Application()
	nt.BrowserTimingHeader()
	nt.NewGoroutine()
	nt.Header()
	nt.Write(nil)
	nt.WriteHeader(0)
	nt.GetTraceMetadata()
	nt.GetLinkingMetadata()
	nt.IsSampled()
}
