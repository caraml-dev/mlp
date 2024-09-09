package newrelic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	newrelic "github.com/newrelic/go-agent"
)

func TestNoopApp(_ *testing.T) {
	na := NoopApp{}
	_ = na.StartTransaction("test", httptest.NewRecorder(), &http.Request{})
	_ = na.RecordCustomEvent("test", nil)
	_ = na.RecordCustomMetric("test", 0)
	_ = na.WaitForConnection(0)
	na.Shutdown(0)
}

func TestNoopTx(_ *testing.T) {
	nt := NoopTx{
		w: httptest.NewRecorder(),
	}
	_ = nt.End()
	_ = nt.Ignore()
	_ = nt.SetName("test")
	_ = nt.NoticeError(nil)
	_ = nt.AddAttribute("key", "val")
	_ = nt.SetWebRequest(nil)
	_ = nt.SetWebResponse(nil)
	_ = nt.StartSegmentNow()
	_ = nt.CreateDistributedTracePayload()
	_ = nt.AcceptDistributedTracePayload(newrelic.TransportUnknown, nil)
	_ = nt.Application()
	_, _ = nt.BrowserTimingHeader()
	_ = nt.NewGoroutine()
	_ = nt.Header()
	_, _ = nt.Write(nil)
	nt.WriteHeader(0)
	_ = nt.GetTraceMetadata()
	_ = nt.GetLinkingMetadata()
	_ = nt.IsSampled()
}
