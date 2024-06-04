// package webhooks provides a webhook manager that can be used to invoke webhooks for different events.

/*
Usage:
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

2. Define the event to webhook configuration. Optionally, the configuration can be provided in a yaml file
and parsed via the `Config` struct.
In the config file, define the event to webhook mapping for those events as required.
For example, if projects need extra labels from an external source,
we define the webhook config for the `OnProjectCreated` event

```go
webhooks:
  enabled: true
  config:
    OnProjectCreated:
      - url: http://localhost:8081/project_created
        method: POST
        onError: abort
```

3. Call InitializeWebhooks() to get a WebhookManager instance.
   This method will initialize the webhook clients for each event type based on the mapping provided

```go
projectsWebhookManager, err := webhooks.InitializeWebhooks(cfg.Webhooks, service.EventList)
```

4. Call `InvokeWebhooks()` method in the caller code based on the event
*/

package webhooks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/caraml-dev/mlp/api/log"
	"github.com/go-playground/validator/v10"
)

type EventType string
type ServiceType string

type WebhookClient interface {
	Invoke(context.Context, []byte) ([]byte, error)
	IsAsync() bool
	IsFinalResponse() bool
	GetUseDataFrom() string
	GetName() string
}

type simpleWebhookClient struct {
	WebhookConfig
}

func NoOpErrorHandler(err error) error { return err }
func NoOpCallback([]byte) error        { return nil }

func (g *simpleWebhookClient) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	// create http request to webhook
	var content []byte
	err := retry.Do(
		func() error {
			client := http.Client{
				Timeout: time.Duration(*g.Timeout) * time.Second,
			}
			req, err := http.NewRequestWithContext(ctx, g.Method, g.URL, bytes.NewBuffer(payload))
			// TODO: Add option for authentication headers
			if err != nil {
				return err
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Errorf("Error making client request %s", err)
				return err
			}
			defer resp.Body.Close()
			content, err = io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if err := validateWebhookResponse(content); err != nil {
				return err
			}
			// check http status code
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("response status code %d not 200", resp.StatusCode)
			}
			return nil

		}, retry.Attempts(uint(g.NumRetries)), retry.Context(ctx),
	)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (g *simpleWebhookClient) IsAsync() bool {
	return g.Async
}

func (g *simpleWebhookClient) IsFinalResponse() bool {
	return g.FinalResponse
}

func (g *simpleWebhookClient) GetUseDataFrom() string {
	return g.UseDataFrom
}

func (g *simpleWebhookClient) GetName() string {
	return g.Name
}

func validateWebhookConfig(webhookConfig *WebhookConfig) error {
	validate := validator.New()

	err := validate.Struct(webhookConfig)
	if err != nil {
		return fmt.Errorf("failed to validate configuration: %s", err)
	}
	if webhookConfig.NumRetries < 0 {
		return fmt.Errorf("numRetries must be a non-negative integer")
	}
	return nil
}

func setDefaults(webhookConfig *WebhookConfig) {
	if webhookConfig.Method == "" {
		webhookConfig.Method = http.MethodPost // Default to POST, TODO: decide if GET is allowed
	}
	if webhookConfig.Timeout == nil {
		def := 10
		webhookConfig.Timeout = &def
	}
}
