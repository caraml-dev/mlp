package webhooks

import "fmt"

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
