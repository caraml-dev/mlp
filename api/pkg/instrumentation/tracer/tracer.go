package tracer

import (
	"github.com/opentracing/opentracing-go"
	"io"

	jaegerconfig "github.com/uber/jaeger-client-go/config"
)

// Config stores tracer configuration
type Config struct {
	// Enabled to enable global tracer set it to true
	Enabled bool

	// ServiceName is service name to be traces
	ServiceName string

	// AgentAddress is jaeger-agent's HTTP sampling server (e.g. localhost:6565)
	AgentAddress string

	// SamplerType is type of sampler: const, probabilistic, rateLimiting, or remote
	SamplerType string

	// SamplerParam is a value passed to the sampler.
	// Valid values for Param field are:
	// - for "const" sampler, 0 or 1 for always false/true respectively
	// - for "probabilistic" sampler, a probability between 0 and 1
	// - for "rateLimiting" sampler, the number of spans per second
	// - for "remote" sampler, param is the same as for "probabilistic"
	SamplerParam float64

	// Tags
	Tags map[string]string
}

type noopCloser struct{}

// Close current trace
func (*noopCloser) Close() error { return nil }

var closer io.Closer = &noopCloser{}

// InitTracer creates a new Jaeger tracer, and sets it as global tracer.
func InitTracer(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}

	tags := make([]opentracing.Tag, 0)
	for k, v := range cfg.Tags {
		tags = append(tags, opentracing.Tag{Key: k, Value: v})
	}

	options := jaegerconfig.Configuration{
		ServiceName: cfg.ServiceName,
		Reporter: &jaegerconfig.ReporterConfig{
			LocalAgentHostPort: cfg.AgentAddress,
		},
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  cfg.SamplerType,
			Param: cfg.SamplerParam,
		},
		Tags: tags,
	}

	tracer, c, err := options.NewTracer()
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)
	closer = c
	return nil
}

// Close flushes buffers before shutdown.
func Close() error {
	return closer.Close()
}
