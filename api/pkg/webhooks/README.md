# Webhooks

- The package is meant to be used across caraml components (e.g. merlin, turing, mlp) to call webhooks when specific events occur.
- The package contains the webhook client implementation and abstracts the logic from the user. It provides some helper functions for the user to call in their code when specific events occur.
- The payload to the webhook server and the response can be arbitrary, and it is up to the user to choose what payload to send to the webhook server(s), but only 1 response will be used in the callback

### How to use?

1. In the caller package (eg, mlp, merlin), define the list of events that requires webhooks. For example:

```go
const (
	ProjectCreatedEvent wh.EventType = "OnProjectCreated"
	ProjectUpdatedEvent wh.EventType = "OnProjectUpdated"
)

var EventList = []wh.EventType{
	ProjectCreatedEvent,
	ProjectUpdatedEvent,
}
```

2. Define the event to webhook configuration. Optionally, the configuration can be provided in a yaml file and parsed via the `Config` struct. In the config file, define the event to webhook mapping for those events as required. For example, if projects need extra labels from an external source, we define the webhook config for the `OnProjectCreated` event

```yaml
webhooks:
  enabled: true
  config:
    OnProjectCreated:
      - url: http://localhost:8081/project_created
        method: POST
        finalResponse: true
        name: webhook1
```

3. Call InitializeWebhooks() to get a WebhookManager instance.
   This method will initialize the webhook clients for each event type based on the mapping provided

```go
projectsWebhookManager, err := webhooks.InitializeWebhooks(cfg.Webhooks, service.EventList)
```

4. Call

```go
InvokeWebhooks(context.Context, EventType, payload interface{}, onSuccess func([]byte) error, onError func(error) error) error
```

method in the caller code based on the event.

### Single Webhook Configuration

```yaml
webhooks:
  enabled: true
  config:
    OnProjectCreated:
      - name: webhook1
        url: http://webhook1
        method: POST
        finalResponse: true
```

- This configuration is the most straight forward. It configures 1 webhook client to be called when the `OnProjectCreated` event happens.
- The payload to the webhook is the json payload of the `payload` argument passed to `InvokeWebhooks`.
- The response from this webhook is used as the final response to the callback passed to the `onSuccess` argument.

### Multiple Webhooks use case

- The library supports multiple webhooks per event to a certain extent.

#### Use case 1

- sync and async webhook
- This can be specified by:

```yaml
webhooks:
  enabled: true
  config:
    OnProjectCreated:
      - name: webhook1
        url: http://webhook1
        method: POST
        finalResponse: true
      - name: webhook2
        url: http://webhook2
        method: POST
        async: true
```

- The async webhook2 will be called only after webhook1 completes.
- If there are multiple sync and async webhooks, the async webhooks will be called only after all sync webhooks have completed.

#### Use case 2

- 3 sync clients, where the response of the first webhook is used as the payload for the second webhook.
- This can be specified by:

```yaml
webhooks:
  enabled: true
  config:
    OnProjectCreated:
      - url: http://webhook1
        method: POST
        finalResponse: true
        name: webhook1
      - url: http://webhook2
        method: POST
        useDataFrom: webhook1 # <-- specify to use data from webhook1
        name: webhook2
      - url: http://webhook3
        method: POST
        name: webhook3
```

- The order of webhook matters, webhook1 will be called before webhook2. If webhook2 is defined before webhook1 but uses the response from webhook1, there will be a validation error on initialization.
- Since `useDataFrom` for webhook1 is not set, webhook1 uses the original payload passed to `InvokeWebhooks` function.
- webhook2 will use the response from webhook1 as its payload. The response from webhook2 is not used.
- webhook3 will use the same payload as webhook1, but will only be called after webhook2
- Here, the finalResponse is set to true for webhook1. This means that the response from webhook1 will be passed as an argument to the `onSuccess` function

### Error Handling

- For synchronous webhooks, all webhooks must be successful before the `onSuccess` handler is called. This means that the caller of this package
  only needs to consider how to handle the successful response.
- In the event any sync webhooks fail, the `onError` handler is called
- For webhooks that do not need to succeed (for whatever reason), pass them as async webhooks.
