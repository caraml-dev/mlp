package newrelic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitNewRelic(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			"disabled",
			Config{
				Enabled: false,
			},
			false,
		},
		{
			"no license",
			Config{
				Enabled: true,
			},
			false,
		},
		{
			"false license",
			Config{
				Enabled: true,
				AppName: "dummy",
				License: "dummy",
			},
			true,
		},
		{
			"dummy license",
			Config{
				Enabled: true,
				AppName: "dummy",
				License: "1234567890123456789012345678901234567890",
			},
			false,
		},
		{
			"with labels",
			Config{
				Enabled: true,
				AppName: "dummy",
				License: "1234567890123456789012345678901234567890",
				Labels: map[string]interface{}{
					"foo": "bar",
				},
			},
			false,
		},
		{
			"with ignore status codes",
			Config{
				Enabled: true,
				AppName: "dummy",
				License: "1234567890123456789012345678901234567890",
				Labels: map[string]interface{}{
					"foo": "bar",
				},
				IgnoreStatusCodes: []int{400, 404},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newRelicApp = &NoopApp{}

			if err := InitNewRelic(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("InitNewRelic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapHandleFunc(t *testing.T) {
	pattern, handler := WrapHandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	assert.Equal(t, "/ping", pattern)
	assert.NotNil(t, handler)
}

type pingHandler struct{}

func (h *pingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ping" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte("pong"))
}

func TestWrapHandle(t *testing.T) {
	pattern, handler := WrapHandle("/ping", new(pingHandler))
	assert.Equal(t, "/ping", pattern)
	assert.NotNil(t, handler)

	w := httptest.NewRecorder()
	assert.NotNil(t, w)

	r, err := http.NewRequest("GET", "http://localhost:8080/ping", nil)
	assert.Nil(t, err)
	assert.NotNil(t, r)

	handler.ServeHTTP(w, r)
	assert.Equal(t, w.Body.String(), "pong")
}

func TestStartTransaction(t *testing.T) {
	tx := StartTransaction("", nil, nil)
	assert.NotNil(t, tx)
}

func TestRecordCustomEvent(t *testing.T) {
	err := RecordCustomEvent("custom_event", nil)
	assert.Nil(t, err)
}

func TestRecordCustomMetric(t *testing.T) {
	err := RecordCustomMetric("custom_metric", 0)
	assert.Nil(t, err)
}

func TestShutdown(t *testing.T) {
	Shutdown(100)
}
