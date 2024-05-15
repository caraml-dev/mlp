package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/caraml-dev/mlp/api/log"
)

type EventType string
type ServiceType string

const (
	onErrorDefault = "ignore"
	onErrorAbort   = "abort"
)

type WebhookManager interface {
	InvokeWebhooks(context.Context, EventType, interface{}, func([]byte) error, func(error) error) error
}

type webhookManager struct {
	webhookClients map[EventType][]WebhookClient
}

type Config struct {
	Enabled bool
	Config  map[EventType][]WebhookConfig `validate:"required_if=Enabled True"`
}

// InvokeWebhooks iterates through the webhooks for a given event and invokes them.
// Sync webhooks are called first, and only after all of them succeed, the async webhooks are called.
// Sync webhooks are called in the order that they are defined. The call order of async webhooks are
// is not guaranteed.
// If any of the sync clients are set to abort, the whole chain aborts as long as 1 sync request returns error.
// onSuccess and onError are callbacks that are called after all webhooks are invoked.
// For sync clients, the payload into a subsequent webhook is the result of the previous webhook
// as long as the call is successful and the response is a valid json object.
// The webhook's response can either be empty or a valid json response.
func (w *webhookManager) InvokeWebhooks(ctx context.Context, event EventType, p interface{}, onSuccess func([]byte) error, onError func(error) error) error {
	var asyncClients []WebhookClient
	var syncClients []WebhookClient
	whc, ok := w.webhookClients[event]
	if !ok {
		return fmt.Errorf("Could not find event %s", event)
	}
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}
	for _, client := range whc {
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
		if err == nil || len(p) > 0 {
			// only update tmpPayload if no error and
			// payload len is not 0
			tmpPayload = p
		}
		if err != nil && client.AbortOnFail() {
			return onError(err)
		}
	}
	for _, client := range asyncClients {
		// TODO: figure out how to handle errors, especially if each async func is invoked
		// in a separate goroutine
		if err := client.InvokeAsync(ctx, payload); err != nil {
			return onError(err)
		}
	}
	// tmpPayload here is the last response.
	if err := onSuccess(tmpPayload); err != nil {
		return nil
	}

	return nil
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
	URL         string `yaml:"url" validate:"required,url"`
	Method      string `yaml:"method"`
	AuthEnabled bool   `yaml:"authEnabled"`
	AuthToken   string `yaml:"authToken" validate:"required_if=AuthEnabled True"`
	OnError     string `yaml:"onError"`
	Async       bool   `yaml:"async"`
}

func NoOpErrorHandler(err error) error { return err }

func (g *SimpleWebhookClient) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	// create http request to webhook
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, g.Method, g.URL, bytes.NewBuffer(payload))
	// TODO: Add option for authentication headers
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error making client request %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := validateWebhookResponse(content); err != nil {
		return nil, err
	}
	// check http status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code %d not 200", resp.StatusCode)
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

func parseWebhookConfig(eventList []EventType, webhookConfigMap map[EventType][]WebhookConfig) (WebhookManager, error) {
	eventToWHMap := make(map[EventType][]WebhookClient)
	for _, eventType := range eventList {
		if webhookConfigList, ok := webhookConfigMap[eventType]; ok {
			var result []WebhookClient
			for _, webhookConfig := range webhookConfigList {
				if err := validateWebhookConfig(&webhookConfig); err != nil {
					return nil, err
				}
				result = append(result, &SimpleWebhookClient{
					WebhookConfig: webhookConfig,
				})
			}
			eventToWHMap[eventType] = result
		}
	}
	return &webhookManager{webhookClients: eventToWHMap}, nil
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

// validateWebhookResponse ensures that the response from a webhook is either
// a valid json object or empty str
// This is only required for synchronous webhooks, where the webhook's response
// may be used as input to the next webhook or in the callback function
func validateWebhookResponse(content []byte) error {
	if len(content) == 0 {
		return nil
	}
	if json.Valid(content) {
		return nil
	}
	return fmt.Errorf("webhook response is not a valid json object and not empty")
}

func InitializeWebhooks(cfg *Config, eventList []EventType) (WebhookManager, error) {
	if cfg == nil || cfg.Enabled {
		return nil, nil
	}
	wi, err := parseWebhookConfig(eventList, cfg.Config)
	if err != nil {
		return nil, err
	}
	return wi, nil

}
