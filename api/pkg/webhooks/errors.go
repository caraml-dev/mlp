package webhooks

import (
	"encoding/json"
	"fmt"
	"strings"
)

type WebhookError struct {
	msg string
	err error
}

func NewWebhookError(err error) *WebhookError {
	return &WebhookError{
		msg: "error invoking webhook",
		err: err,
	}
}

// Implement errors.Error method
func (e *WebhookError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e *WebhookError) Unwrap() error {
	return e.err
}

func (e *WebhookError) Is(target error) bool {
	_, ok := target.(*WebhookError)
	return ok
}

// ResponseError is returned when a webhook answers with a 4xx status code.
// It keeps the status code and the webhook's own error message so callers
// can turn it into an error the user can act on.
type ResponseError struct {
	WebhookName string
	StatusCode  int
	Message     string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("response from %s, status code %d not 200, err: %s", e.WebhookName, e.StatusCode, e.Message)
}

// webhookErrorMessage extracts the "error" field from a JSON error body such as
// {"error":"team not found"}, and falls back to the raw body.
func webhookErrorMessage(content []byte) string {
	var body struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(content, &body); err == nil && body.Error != "" {
		return body.Error
	}
	if msg := strings.TrimSpace(string(content)); msg != "" {
		return msg
	}
	return "no error message returned"
}
