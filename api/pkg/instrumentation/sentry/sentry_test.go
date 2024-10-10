package sentry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitSentry(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"disabled",
			args{
				Config{
					Enabled: false,
				},
			},
			false,
		},
		{
			"dummy dsn",
			args{
				Config{
					Enabled: true,
					DSN:     "1234567890",
				},
			},
			true,
		},
		{
			"empty dsn",
			args{
				Config{
					Enabled: true,
					DSN:     "",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitSentry(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("InitSentry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSentry(t *testing.T) {
	sentry := Sentry()
	assert.NotNil(t, sentry)

	panicHandler := RecoveryHandler(func(_ http.ResponseWriter, _ *http.Request) {
		panic("at the disco")
	})
	assert.NotNil(t, panicHandler)

	mux := http.NewServeMux()
	mux.Handle("/panic", http.HandlerFunc(panicHandler))

	r, err := http.NewRequest("GET", "http://localhost:8080/panic", nil)
	if err != nil {
		t.Fatalf("Error building test request: %s", err)
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	Close()
}
