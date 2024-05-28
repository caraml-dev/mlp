package webhooks

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking onSuccess and onError functions
var onSuccess = func(t *testing.T, expected []byte) func(response []byte) error {
	return func(input []byte) error {
		if diff := cmp.Diff(input, expected); diff != "" {
			t.Errorf("unexpected response (-got +want):\n%s", diff)
		}
		return nil
	}
}
var onError = func(err error) error {
	fmt.Printf("err: %s", err.Error())
	return err
}

type testPayload struct {
	Data string `json:"data"`
}

var testPayloadData = testPayload{Data: "abc"}

func TestInitializeWebhooks(t *testing.T) {
	tests := []struct {
		name          string
		cfg           *Config
		eventList     []EventType
		expectedError bool
	}{
		{
			name:          "Config is nil",
			cfg:           nil,
			eventList:     nil,
			expectedError: false,
		},
		{
			name: "Config is disabled",
			cfg: &Config{
				Enabled: false,
				Config:  nil,
			},
			eventList:     nil,
			expectedError: false,
		},
		{
			name: "Config enabled with valid webhooks",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:          "webhook1",
							URL:           "http://example.com",
							Method:        "POST",
							AuthEnabled:   false,
							FinalResponse: true,
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: false,
		},
		{
			name: "Config enabled with invalid webhook (missing name)",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:          "",
							URL:           "http://example.com",
							Method:        "POST",
							AuthEnabled:   false,
							FinalResponse: true,
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: true,
		},
		{
			name: "Config enabled with invalid webhook (missing URL)",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:          "webhook1",
							URL:           "",
							Method:        "POST",
							AuthEnabled:   false,
							FinalResponse: true,
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: true,
		},
		{
			name: "Config enabled with invalid webhook (auth enabled but no token)",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:          "webhook1",
							URL:           "http://example.com",
							Method:        "POST",
							AuthEnabled:   true,
							AuthToken:     "",
							FinalResponse: true,
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: true,
		},
		{
			name: "Config enabled with multiple webhooks",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:          "webhook1",
							URL:           "http://example.com",
							Method:        "POST",
							FinalResponse: true,
						},
						{
							Name:   "webhook2",
							URL:    "http://example.com",
							Method: "POST",
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: false,
		},
		{
			name: "Config enabled with multiple webhooks fail (both set FinalResponse)",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:          "webhook1",
							URL:           "http://example.com",
							Method:        "POST",
							FinalResponse: true,
						},
						{
							Name:          "webhook2",
							URL:           "http://example.com",
							Method:        "POST",
							FinalResponse: true,
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: true,
		},
		{
			name: "Config enabled with multiple webhooks fail (duplicate names)",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:          "webhook1",
							URL:           "http://example.com",
							Method:        "POST",
							FinalResponse: true,
						},
						{
							Name:   "webhook1",
							URL:    "http://example.com",
							Method: "POST",
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: true,
		},
		{
			name: "Config enabled with multiple webhooks fail (invalid reference)",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:          "webhook1",
							URL:           "http://example.com",
							Method:        "POST",
							FinalResponse: true,
						},
						{
							Name:        "webhook2",
							URL:         "http://example.com",
							Method:      "POST",
							UseDataFrom: "webhook_non_existent",
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: true,
		},
		{
			name: "Config enabled with multiple webhooks fail (incorrect order)",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:        "webhook2",
							URL:         "http://example.com",
							Method:      "POST",
							UseDataFrom: "webhook1",
						},
						{
							Name:          "webhook1",
							URL:           "http://example.com",
							Method:        "POST",
							FinalResponse: true,
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: true,
		},
		{
			name: "Config enabled with multiple webhooks with sync and async, fail with duplicate names",
			cfg: &Config{
				Enabled: true,
				Config: map[EventType][]WebhookConfig{
					"event1": {
						{
							Name:        "webhook2",
							URL:         "http://example.com",
							Method:      "POST",
							UseDataFrom: "webhook1",
							Async:       true,
						},
						{
							Name:          "webhook1",
							URL:           "http://example.com",
							Method:        "POST",
							FinalResponse: true,
						},
						{
							Name:   "webhook2",
							URL:    "http://example.com",
							Method: "POST",
							Async:  true,
						},
					},
				},
			},
			eventList:     []EventType{"event1"},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := InitializeWebhooks(test.cfg, test.eventList)
			if err != nil && !test.expectedError || err == nil && test.expectedError {
				t.Errorf("expected error: %v, got: %v", test.expectedError, err)
			}
		})
	}
}
func TestInvokeWebhooksSimple(t *testing.T) {
	// Setup mock WebhookClient with specific behaviors
	response := []byte(`{"data":"abc"}`)
	mockClient := &MockWebhookClient{}
	// Expect Invoke to return response from client and no error
	mockClient.On("Invoke", mock.Anything, mock.Anything).Return(response, nil).Once()
	mockClient.On("IsFinalResponse").Return(true)
	mockClient.On("GetUseDataFrom").Return("")
	mockClient.On("GetName").Return("webhook1")

	// Setup WebhookManager with the mock client
	webhookManager := &SimpleWebhookManager{
		WebhookClients: map[EventType]map[WebhookType][]WebhookClient{
			"validEvent": {Sync: {mockClient}},
		},
	}

	// Execution
	err := webhookManager.InvokeWebhooks(
		context.Background(),
		"validEvent",
		&testPayloadData,
		onSuccess(t, response),
		onError,
	)

	// Assertion
	assert.NoError(t, err)           // Expect no error
	mockClient.AssertExpectations(t) // Verify that expectations on the mock client were met
}

func TestInvokeMultipleSyncWebhooks(t *testing.T) {
	// Setup mock WebhookClient with specific behaviors
	webhook1Result := []byte(`{"result": "xyz"}`)
	response := []byte(`{"data":"abc"}`)
	mockClient := &MockWebhookClient{}
	// Expect Invoke to return response from client and no error
	mockClient.On("Invoke", mock.Anything, mock.Anything).Return(webhook1Result, nil)
	mockClient.On("GetName").Return("webhook1")
	mockClient.On("IsFinalResponse").Return(false)
	mockClient.On("GetUseDataFrom").Return("")

	mockClient2 := &MockWebhookClient{}
	mockClient2.On("Invoke", mock.Anything, mock.MatchedBy(func(payload []byte) bool {
		if diff := cmp.Diff(webhook1Result, payload); diff != "" {
			fmt.Println("diff", diff)
			return false
		}
		return true
	})).Return(response, nil)
	mockClient2.On("GetName").Return("webhook2")
	mockClient2.On("GetUseDataFrom").Return("webhook1")
	mockClient2.On("IsFinalResponse").Return(true)
	// Setup WebhookManager with the mock client
	webhookManager := &SimpleWebhookManager{
		WebhookClients: map[EventType]map[WebhookType][]WebhookClient{
			"validEvent": {Sync: {mockClient, mockClient2}},
		},
	}
	// Execution
	err := webhookManager.InvokeWebhooks(
		context.Background(),
		"validEvent",
		&testPayloadData,
		onSuccess(t, response),
		onError,
	)

	// Assertion
	assert.NoError(t, err)
	mockClient.AssertExpectations(t) // Verify that expectations on the mock client were met
}
