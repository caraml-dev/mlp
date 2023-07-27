package enforcer

// Builder builder of enforcer.Enforcer
type Builder struct {
	ketoRemoteRead  string
	ketoRemoteWrite string
	cacheConfig     *CacheConfig
}

const (
	// DefaultKetoRemoteRead default Keto remote read endpoint
	DefaultKetoRemoteRead = "http://localhost:4466"
	// DefaultKetoRemoteWrite default Keto remote write endpoint
	DefaultKetoRemoteWrite = "http://localhost:4467"
)

// NewEnforcerBuilder create new enforcer builder with all default parameters
func NewEnforcerBuilder() *Builder {
	return &Builder{
		ketoRemoteRead:  DefaultKetoRemoteRead,
		ketoRemoteWrite: DefaultKetoRemoteWrite,
	}
}

// KetoEndpoints set Keto remote read and write endpoint
func (b *Builder) KetoEndpoints(ketoRemoteRead string, ketoRemoteWrite string) *Builder {
	b.ketoRemoteRead = ketoRemoteRead
	b.ketoRemoteWrite = ketoRemoteWrite
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
	return newEnforcer(b.ketoRemoteRead, b.ketoRemoteWrite, b.cacheConfig)
}
