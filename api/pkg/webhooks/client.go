package webhooks

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

type EventType string
type ServiceType string

const (
	ProjectServiceType ServiceType = "project"
)

const (
	ProjectCreatedEventType EventType = "OnProjectCreated"
)

const (
	onErrorDefault = "ignore"
	onErrorAbort   = "abort"
)

// Define serviceEventMapping here
var serviceEventMapping = map[ServiceType][]EventType{
	ProjectServiceType: {ProjectCreatedEventType},
}

type WebhookManagerI interface {
	InvokeWebhooks(context.Context, []byte, func([]byte, interface{}) error, func(error) error) error
	GetObject() interface{}
}

type WebhookManager struct {
	WebhookClients []WebhookClient
	Object         interface{}
}

// InvokeWebhooks iterates through sync clients and async clients
// For sync clients, preserve order. If any of the sync clients are set to abort, the whole chain aborts as long as 1 sync request returns error
// onSuccess and onError are callbacks that are called after all webhooks are invoked.
// For sync clients, the payload into a subsequent webhook is the result of the previous webhook
func (w *WebhookManager) InvokeWebhooks(ctx context.Context, payload []byte, onSuccess func([]byte, interface{}) error, onError func(error) error) error {
	var asyncClients []WebhookClient
	var syncClients []WebhookClient
	for _, client := range w.WebhookClients {
		if client.IsAsync() {
			asyncClients = append(asyncClients, client)
		} else {
			syncClients = append(syncClients, client)
		}
	}
	// create copy of original payload
	tmpPayload := make([]byte, len(payload))
	copy(tmpPayload, payload)
	for _, client := range syncClients {
		p, err := client.Invoke(ctx, tmpPayload)
		if err != nil && client.AbortOnFail() {
			return onError(err)
		}
		tmpPayload = p
	}
	for _, client := range asyncClients {
		// TODO: figure out how to handle errors, especially if each async func is invoked
		// in a separate goroutine
		if err := client.InvokeAsync(ctx, payload); err != nil {
			return onError(err)
		}
	}
	// tmpPayload here is the last response.
	if err := onSuccess(tmpPayload, w.GetObject()); err != nil {
		return nil
	}

	return nil
}

func (w *WebhookManager) GetObject() interface{} {
	return w.Object
}

type WebhookClient interface {
	Invoke(context.Context, []byte) ([]byte, error)
	InvokeAsync(context.Context, []byte) error
	IsAsync() bool
	AbortOnFail() bool
}

type SimpleWebhookClient struct {
	WebhookConfig
}

type WebhookConfig struct {
	URL         string `yaml:"url"`
	Method      string `yaml:"method"`
	AuthEnabled bool   `yaml:"authEnabled"`
	AuthToken   string `yaml:"authToken"`
	OnError     string `yaml:"onError"`
	Async       bool   `yaml:"async"`
}

func NoOpErrorHandler(_ error) {}

func (g *SimpleWebhookClient) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	// create http request to webhook
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, g.Method, g.URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var content []byte
	if _, err := resp.Body.Read(content); err != nil {
		return nil, err
	}
	return content, nil
}

func (g *SimpleWebhookClient) InvokeAsync(ctx context.Context, payload []byte) error {
	go func() {
		if _, err := g.Invoke(ctx, payload); err != nil {
			return
		}
	}()
	return nil
}

func (g *SimpleWebhookClient) IsAsync() bool {
	return g.Async
}

func (g *SimpleWebhookClient) AbortOnFail() bool {
	return g.OnError == onErrorAbort
}

func ParseWebhookConfig(serviceType ServiceType, webhookConfigMap map[EventType][]WebhookConfig) ([]WebhookClient, error) {
	var result []WebhookClient
	availableEvents, ok := serviceEventMapping[serviceType]
	if !ok {
		// invalid serviceType passed in, return empty slice
		return nil, nil
	}
	for _, eventType := range availableEvents {
		if webhookConfigList, ok := webhookConfigMap[eventType]; ok {
			for _, webhookConfig := range webhookConfigList {
				if err := validateWebhookConfig(&webhookConfig); err != nil {
					return nil, err
				}
				result = append(result, &SimpleWebhookClient{
					WebhookConfig: webhookConfig,
				})
			}
		}
	}

	return result, nil

}

func validateWebhookConfig(webhookConfig *WebhookConfig) error {
	if webhookConfig.URL == "" {
		return fmt.Errorf("missing webhook URL")
	}
	if webhookConfig.Method == "" {
		webhookConfig.Method = http.MethodGet
	}
	if webhookConfig.AuthEnabled && webhookConfig.AuthToken == "" {
		return fmt.Errorf("missing webhook auth token")
	}
	if webhookConfig.OnError == "" {
		webhookConfig.OnError = onErrorDefault
	}
	return nil
}
