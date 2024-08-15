package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
)

type WebhookManager interface {
	InvokeWebhooks(
		context.Context,
		EventType,
		interface{},
		func(payload []byte) error,
		func(error) error,
	) error

	IsEventConfigured(EventType) bool
}

type SimpleWebhookManager struct {
	SyncClients  map[EventType][]WebhookClient
	AsyncClients map[EventType][]WebhookClient
}

// IsEventConfigured checks if the event is configured in the webhook manager
// Use this method before calling InvokeWebhooks if it is optional to set webhooks for an event
func (w *SimpleWebhookManager) IsEventConfigured(event EventType) bool {
	_, ok := w.SyncClients[event]
	_, ok1 := w.AsyncClients[event]
	return ok || ok1
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
func (w *SimpleWebhookManager) InvokeWebhooks(
	ctx context.Context,
	event EventType,
	p interface{},
	onSuccess func([]byte) error,
	onError func(error) error,
) error {
	finalResponse := make([]byte, 0)
	syncClients, ok := w.SyncClients[event]
	asyncClients, ok1 := w.AsyncClients[event]
	if !ok && !ok1 {
		return fmt.Errorf("Could not find event %s", event)
	}
	originalPayload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// Mapping to store response from different webhooks
	responsePayloadLookup := make(map[string][]byte)

	for _, client := range syncClients {
		var tmpPayload []byte
		if client.GetUseDataFrom() == "" {
			tmpPayload = originalPayload
		} else {
			tmpPayload, ok = responsePayloadLookup[client.GetUseDataFrom()]
			if !ok {
				// NOTE: This should never happen!
				return fmt.Errorf(
					"webhook name %s not found, this could be because of an error in a previous webhook that this webhook depends on",
					client.GetUseDataFrom(),
				)
			}
		}
		p, err := client.Invoke(ctx, tmpPayload)
		if err != nil {
			return onError(err)
		}
		responsePayloadLookup[client.GetName()] = p
		if client.IsFinalResponse() {
			finalResponse = p
		}
	}
	for _, client := range asyncClients {
		go func(client WebhookClient) {
			// Ignore the response from async webhooks
			if _, err := client.Invoke(context.Background(), originalPayload); err != nil {
				return
			}
		}(client)
	}
	if err := onSuccess(finalResponse); err != nil {
		// If the callback fails, return the error
		return err
	}

	return nil
}

func parseAndValidateConfig(
	eventList []EventType,
	webhookConfigMap map[EventType][]WebhookConfig,
) (WebhookManager, error) {
	syncClientMap := make(map[EventType][]WebhookClient)
	asyncClientMap := make(map[EventType][]WebhookClient)
	for _, eventType := range eventList {
		webhookConfigList, ok := webhookConfigMap[eventType]
		if !ok {
			continue
		}
		syncClients := make([]WebhookClient, 0)
		asyncClients := make([]WebhookClient, 0)
		for _, webhookConfig := range webhookConfigList {
			if err := validateWebhookConfig(&webhookConfig); err != nil {
				return nil, err
			}
			setDefaults(&webhookConfig)
			client := &simpleWebhookClient{
				WebhookConfig: webhookConfig,
			}
			if !client.IsAsync() {
				syncClients = append(syncClients, client)
			} else {
				asyncClients = append(asyncClients, client)
			}
		}
		allClients := append(syncClients, asyncClients...)
		if err := validateClients(allClients); err != nil {
			return nil, err
		}
		syncClientMap[eventType] = syncClients
		asyncClientMap[eventType] = asyncClients
	}
	return &SimpleWebhookManager{AsyncClients: asyncClientMap, SyncClients: syncClientMap}, nil

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

func validateClients(webhookClients []WebhookClient) error {
	// ensure that only 1 sync client has finalResponse set to true
	isFinalResponseSet := false
	// Check for duplicate webhook names
	webhookNames := make(map[string]int)
	for idx, client := range webhookClients {
		if client.IsFinalResponse() {
			if isFinalResponseSet {
				return fmt.Errorf("only 1 sync client can have finalResponse set to true")
			}
			isFinalResponseSet = true
		}
		if _, ok := webhookNames[client.GetName()]; ok {
			return fmt.Errorf("duplicate webhook name")
		}
		webhookNames[client.GetName()] = idx
		// Ensure that if a client uses the response from another client, that client exists
		// If a client uses the response from another client, it must be defined before it
		if client.GetUseDataFrom() == "" {
			// If the client does not use data from another webhook,
			// then we can skip the rest of the checks
			// since the payload used will be the user's payload
			continue
		}
		useIdx, ok := webhookNames[client.GetUseDataFrom()]
		if !ok {
			return fmt.Errorf("webhook name %s not found", client.GetUseDataFrom())
		}
		if useIdx > idx {
			return fmt.Errorf(
				"webhook name %s must be defined before %s",
				client.GetUseDataFrom(),
				client.GetName(),
			)
		}
	}
	if !isFinalResponseSet {
		return fmt.Errorf("at least 1 sync client must have finalResponse set to true")
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
