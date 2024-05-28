package webhooks

type WebhookType string

const (
	Async WebhookType = "async"
	Sync  WebhookType = "sync"
)

// Config is a helper struct to define the webhook config in a configuration file
type Config struct {
	Enabled bool
	Config  map[EventType][]WebhookConfig `validate:"required_if=Enabled True"`
}

// WebhookConfig struct is the configuration for each webhook to be called
type WebhookConfig struct {
	Name        string `yaml:"name"        validate:"required"`
	URL         string `yaml:"url"         validate:"required,url"`
	Method      string `yaml:"method"`
	AuthEnabled bool   `yaml:"authEnabled"`
	AuthToken   string `yaml:"authToken"   validate:"required_if=AuthEnabled True"`
	OnError     string `yaml:"onError"`
	Async       bool   `yaml:"async"`
	NumRetries  int    `yaml:"numRetries"`
	Timeout     *int   `yaml:"timeout"`
	// UseDataFrom is the name of the webhook whose response will be used as input to this webhook
	UseDataFrom string `yaml:"useDataFrom"`

	// FinalResponse can be set to use the response from this webhook to the onSuccess callback function
	FinalResponse bool `yaml:"finalResponse"`
}
