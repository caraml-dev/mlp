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

// webhookErrorMessage returns the "error" field of a JSON error body.
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
