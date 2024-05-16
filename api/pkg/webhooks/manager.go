package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhookManager interface {
	InvokeWebhooks(context.Context, EventType, interface{}, func([]byte) error, func(error) error) error
}

type webhookManager struct {
	webhookClients map[EventType][]WebhookClient
}

// InvokeWebhooks iterates through the webhooks for a given event and invokes them.
// Sync webhooks are called first, and only after all of them succeed, the async webhooks are called.
// Sync webhooks are called in the order that they are defined. The call order of async webhooks are
// is not guaranteed.
// If any of the sync clients are set to abort, the whole chain aborts as long as 1 sync request returns error.
// onSuccess and onError are callbacks that are called after all webhooks are invoked.
// For sync clients, the payload can be either the original input payload, or the response from another sync webhook.
// This can be specified in the UseDataFrom field
// For async clients, the payload is only the original input payload.
// Only one webhook's response can be used as the finalResponse
func (w *webhookManager) InvokeWebhooks(ctx context.Context, event EventType, p interface{}, onSuccess func([]byte) error, onError func(error) error) error {
	var asyncClients []WebhookClient
	var syncClients []WebhookClient
	var finalResponse []byte
	whc, ok := w.webhookClients[event]
	if !ok {
		return fmt.Errorf("Could not find event %s", event)
	}
	originalPayload, err := json.Marshal(p)
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

	// Mapping to store response from different webhooks
	responsePayloadLookup := make(map[string][]byte)

	for _, client := range syncClients {
		var tmpPayload []byte
		if client.GetUseDataFrom() == "" {
			tmpPayload = originalPayload
		} else if tmpPayload, ok = responsePayloadLookup[client.GetUseDataFrom()]; !ok {
			// This should only happen if a previous error had an error, but did not abort
			// and the current client is trying to use the response from that client
			return fmt.Errorf("webhook name %s not found, this could be because an error in a downstream webhook", client.GetUseDataFrom())
		}
		p, err := client.Invoke(ctx, tmpPayload)
		if err == nil {
			responsePayloadLookup[client.GetName()] = p
			if client.IsFinalResponse() {
				finalResponse = p
			}
			continue
		}
		// if err is not nil, check if client is set to abort
		if client.AbortOnFail() {
			return onError(err)
		}

	}
	for _, client := range asyncClients {
		// NOTE: Currently, this will never return err since InvokeAsync always returns nil
		if err := client.InvokeAsync(ctx, originalPayload); err != nil {
			return onError(err)
		}
	}
	// tmpPayload here is the last response.
	if err := onSuccess(finalResponse); err != nil {
		// If the callback fails, return the error
		return err
	}

	return nil
}

func parseAndValidateConfig(eventList []EventType, webhookConfigMap map[EventType][]WebhookConfig) (WebhookManager, error) {
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
	for _, webhookClients := range eventToWHMap {
		syncClients := make([]WebhookClient, 0)
		for _, client := range webhookClients {
			if !client.IsAsync() {
				syncClients = append(syncClients, client)
			}
		}
		if err := validateSyncClients(syncClients); err != nil {
			return nil, err
		}
	}
	return &webhookManager{webhookClients: eventToWHMap}, nil
}

func validateWebhookConfig(webhookConfig *WebhookConfig) error {
	if webhookConfig.Name == "" {
			return fmt.Errorf("missing webhook name")
	}
	if webhookConfig.URL == "" {
		return fmt.Errorf("missing webhook URL")
	}
	if webhookConfig.Method == "" {
		webhookConfig.Method = http.MethodPost // Default to POST, TODO: decide if GET is allowed
	}
	if webhookConfig.AuthEnabled && webhookConfig.AuthToken == "" {
		return fmt.Errorf("missing webhook auth token")
	}
	if webhookConfig.OnError == "" {
		webhookConfig.OnError = onErrorAbort
	}
	if webhookConfig.NumRetries < -1 {
		return fmt.Errorf("numRetries must be a positive integer or -1")
	}
	if webhookConfig.Timeout == nil {
		def := 10
		webhookConfig.Timeout = &def
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

func validateSyncClients(webhookClients []WebhookClient) error {
	// ensure that only 1 sync client has finalResponse set to true
	isFinalResponseSet := false
	for _, client := range webhookClients {
		if client.IsFinalResponse() {
			if isFinalResponseSet {
				return fmt.Errorf("only 1 sync client can have finalResponse set to true")
			}
			isFinalResponseSet = true
		}
	}
	if !isFinalResponseSet {
		return fmt.Errorf("at least 1 sync client must have finalResponse set to true")
	}
	// Ensure that all useDataFrom fields exist
	webhookNames := make(map[string]int)
	for idx, client := range webhookClients {
		if _, ok := webhookNames[client.GetName()]; ok {
			return fmt.Errorf("duplicate webhook name")
		}
		webhookNames[client.GetName()] = idx

	}
	// Ensure that webhook order is correct if they have dependencies
	for idx, client := range webhookClients {
		if client.GetUseDataFrom() == "" {
			continue
		}
		useIdx, ok := webhookNames[client.GetUseDataFrom()]
		if !ok {
			return fmt.Errorf("webhook name %s not found", client.GetUseDataFrom())
		}
		if useIdx > idx {
			return fmt.Errorf("webhook name %s must be defined before %s", client.GetUseDataFrom(), client.GetName())
		}
	}

	return nil
}

// InitializeWebhooks is a helper method to initialize a webhook manager based on the eventList
// provided. It returns an error if the configuration is invalid
func InitializeWebhooks(cfg *Config, eventList []EventType) (WebhookManager, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}
	wi, err := parseAndValidateConfig(eventList, cfg.Config)
	if err != nil {
		return nil, err
	}
	return wi, nil

}
