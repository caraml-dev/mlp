package sentry

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	raven "github.com/getsentry/raven-go"
)

// Client is a client to send logs to Sentry.
type Client interface {
	Capture(packet *raven.Packet, captureTags map[string]string) (eventID string, ch chan error)
	CaptureError(err error, tags map[string]string, interfaces ...raven.Interface) string
	Close()
}

var (
	sentry Client = &NoopClient{}
)

// Config stores NewRelic configuration.
type Config struct {
	Enabled bool
	DSN     string
	Labels  map[string]string
}

// InitSentry creates a new Sentry client.
func InitSentry(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}

	client, err := raven.NewWithTags(
		cfg.DSN,
		cfg.Labels,
	)
	if err != nil {
		return err
	}

	sentry = client
	return nil
}

// Sentry returns the singleton Sentry client implementation.
func Sentry() Client {
	return sentry
}

// Close flushes the Client's buffer and releases the associated ressources. The
// Client and all the cloned Clients must not be used afterward.
func Close() {
	sentry.Close()
}

// RecoveryHandler wraps the stdlib net/http Mux.
func RecoveryHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return Recoverer(http.HandlerFunc(handler)).ServeHTTP
}

// Recoverer wraps the stdlib net/http Mux.
func Recoverer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()
				rvalStr := fmt.Sprint(rval)
				var packet *raven.Packet
				if err, ok := rval.(error); ok {
					packet = raven.NewPacket(
						rvalStr,
						raven.NewException(
							errors.New(rvalStr),
							raven.GetOrNewStacktrace(err, 2, 3, nil)),
						raven.NewHttp(r))
				} else {
					packet = raven.NewPacket(
						rvalStr,
						raven.NewException(
							errors.New(rvalStr),
							raven.NewStacktrace(2, 3, nil)),
						raven.NewHttp(r))
				}
				sentry.Capture(packet, nil)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
