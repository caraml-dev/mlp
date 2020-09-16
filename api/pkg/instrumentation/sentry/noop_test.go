package sentry

import (
	"testing"
)

func TestNoopClient(t *testing.T) {
	nc := &NoopClient{}
	nc.Capture(nil, nil)
	nc.CaptureError(nil, nil, nil)
	nc.Close()
}
