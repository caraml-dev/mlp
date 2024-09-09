package sentry

import (
	"testing"
)

func TestNoopClient(_ *testing.T) {
	nc := &NoopClient{}
	nc.Capture(nil, nil)
	nc.CaptureError(nil, nil, nil)
	nc.Close()
}
