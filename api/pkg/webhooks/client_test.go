package webhooks

import (
	"context"
	"fmt"
	"testing"

	"encoding/json"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking onSuccess and onError functions
var onSuccess = func(response []byte) error {
	return nil
}
var onError = func(err error) error {
	fmt.Printf("err: %s", err.Error())
	return err
}

type testPayload struct {
	Data string `json:"data"`
}

type testResult struct {
	Result string `json:"result"`
}

var testPayloadData = testPayload{Data: "abc"}
var testResultData = testResult{Result: "xyz"}

func TestInvokeWebhooksSimple(t *testing.T) {
	// Setup mock WebhookClient with specific behaviors
	mockClient := &MockWebhookClient{}
	// Expect Invoke to return response from client and no error
	mockClient.On("Invoke", mock.Anything, mock.Anything).Return([]byte("{}"), nil).Once()
	mockClient.On("IsAsync").Return(false)
	mockClient.On("AbortOnFail").Return(false)

	// Setup WebhookManager with the mock client
	webhookManager := &webhookManager{
		webhookClients: map[EventType][]WebhookClient{
			"validEvent": {mockClient},
		},
	}

	// Execution
	err := webhookManager.InvokeWebhooks(context.Background(), "validEvent", &testPayloadData, onSuccess, onError)

	// Assertion
	assert.NoError(t, err)           // Expect no error
	mockClient.AssertExpectations(t) // Verify that expectations on the mock client were met
}

func TestInvokeMultipleSyncWebhooks(t *testing.T) {
	// Setup mock WebhookClient with specific behaviors
	mockClient := &MockWebhookClient{}
	// Expect Invoke to return response from client and no error
	mockClient.On("Invoke", mock.Anything, mock.Anything).Return([]byte(`{"result": "xyz"}`), nil)
	mockClient.On("IsAsync").Return(false)

	mockClient2 := &MockWebhookClient{}
	mockClient2.On("Invoke", mock.Anything, mock.MatchedBy(func(payload []byte) bool {
		// check if payload matches testPayloadData
		var tmp testResult
		if err := json.Unmarshal(payload, &tmp); err != nil {
			return false
		}
		if diff := cmp.Diff(testResultData, tmp); diff != "" {
			fmt.Println("diff", diff)
			return false
		}
		return true
	})).Return([]byte("{}"), nil)
	mockClient2.On("IsAsync").Return(false).Once()
	// Setup WebhookManager with the mock client
	webhookManager := &webhookManager{
		webhookClients: map[EventType][]WebhookClient{
			"validEvent": {mockClient, mockClient2},
		},
	}
	// Execution
	err := webhookManager.InvokeWebhooks(context.Background(), "validEvent", &testPayloadData, onSuccess, onError)

	// Assertion
	assert.NoError(t, err)
	mockClient.AssertExpectations(t) // Verify that expectations on the mock client were met
	// t.Errorf("test")
}
