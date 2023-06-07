package enforcer

import (
	"time"
)

// Builder builder of enforcer.Enforcer
type Builder struct {
	url         string
	product     string
	flavor      Flavor
	timeout     time.Duration
	cacheConfig *CacheConfig
}

const (
	// DefaultURL default Keto server URL
	DefaultURL = "http://localhost:4466"
	// DefaultFlavor default Keto flavor to be used
	DefaultFlavor = FlavorGlob
	// DefaultTimeout maximum call duration to Keto Server before considered as timeout
	DefaultTimeout = 5 * time.Second
)

// NewEnforcerBuilder create new enforcer builder with all default parameters
func NewEnforcerBuilder() *Builder {
	return &Builder{
		url:     DefaultURL,
		flavor:  DefaultFlavor,
		timeout: DefaultTimeout,
	}
}

// Product set product name
func (b *Builder) Product(product string) *Builder {
	b.product = product
	return b
}

// URL set Keto URL
func (b *Builder) URL(url string) *Builder {
	b.url = url
	return b
}

// Flavor set Keto flavor
func (b *Builder) Flavor(flavor Flavor) *Builder {
	b.flavor = flavor
	return b
}

// Timeout set timeout
func (b *Builder) Timeout(timeout time.Duration) *Builder {
	b.timeout = timeout
	return b
}

func (b *Builder) WithCaching(keyExpirySeconds int, cacheCleanUpIntervalSeconds int) *Builder {
	b.cacheConfig = &CacheConfig{
		KeyExpirySeconds:            keyExpirySeconds,
		CacheCleanUpIntervalSeconds: cacheCleanUpIntervalSeconds,
	}
	return b
}

// Build build an enforcer.Enforcer instance
func (b *Builder) Build() (Enforcer, error) {
	return newEnforcer(b.url, b.product, b.flavor, b.timeout, b.cacheConfig)
}
