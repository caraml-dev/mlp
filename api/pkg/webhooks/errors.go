package webhooks

import (
	"fmt"

	"github.com/google/uuid"
)

type WebhookError struct {
	msg         string
	err         error
	ReferenceID string
}

// NewWebhookError wraps a failed webhook call. Each error gets its own reference id, which is part of
// the message, so a user can quote it and the platform owner can find the matching log line.
func NewWebhookError(err error) *WebhookError {
	referenceID := uuid.New().String()
	return &WebhookError{
		msg:         fmt.Sprintf("error invoking webhook. Please contact the platform owner (reference id: %s)", referenceID),
		err:         err,
		ReferenceID: referenceID,
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
