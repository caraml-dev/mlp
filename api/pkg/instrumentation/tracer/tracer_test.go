package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitGlobalTracer(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		wantNoop bool
		wantErr  bool
	}{
		{
			"success disabled",
			Config{
				Enabled:      false,
				ServiceName:  "service",
				AgentAddress: ":6565",
				SamplerType:  "probabilistic",
				SamplerParam: 0.01,
				Tags:         map[string]string{"key": "value"},
			},
			true,
			false,
		},
		{
			"success enabled",
			Config{
				Enabled:      true,
				ServiceName:  "service",
				AgentAddress: ":6565",
				SamplerType:  "probabilistic",
				SamplerParam: 0.01,
				Tags:         map[string]string{"key": "value"},
			},
			false,
			false,
		},
		{
			"success no tags",
			Config{
				Enabled:      true,
				ServiceName:  "service",
				AgentAddress: ":6565",
				SamplerType:  "random",
				SamplerParam: 0.01,
			},
			false,
			true,
		},
		{
			"fail no service name",
			Config{
				Enabled:      true,
				ServiceName:  "",
				AgentAddress: ":6565",
				SamplerType:  "probabilistic",
				SamplerParam: 0.01,
				Tags:         map[string]string{"key": "value"},
			},
			false,
			true,
		},
		{
			"fail invalid sampler type",
			Config{
				Enabled:      true,
				ServiceName:  "service",
				AgentAddress: ":6565",
				SamplerType:  "random",
				SamplerParam: 0.01,
				Tags:         map[string]string{"key": "value"},
			},
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitTracer(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("InitTracer() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantNoop {
				assert.IsType(t, &noopCloser{}, closer)
			}
		})
	}
}

func TestClose(t *testing.T) {
	err := Close()
	assert.Nil(t, err)
}

func Test_noopCloser_Close(t *testing.T) {
	n := &noopCloser{}
	err := n.Close()
	assert.Nil(t, err)
}
