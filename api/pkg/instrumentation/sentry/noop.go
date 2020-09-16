package sentry

import raven "github.com/getsentry/raven-go"

// A NoopClient is a trivial, minimum overhead implementation of Client
// for which all operations are no-ops.
type NoopClient struct{}

// Capture implements Client interface.
func (nc *NoopClient) Capture(packet *raven.Packet, captureTags map[string]string) (eventID string, ch chan error) {
	return "", nil
}

// CaptureError implements Client interface.
func (nc *NoopClient) CaptureError(err error, tags map[string]string, interfaces ...raven.Interface) string {
	return ""
}

// Close implements Client interface.
func (nc *NoopClient) Close() {}
