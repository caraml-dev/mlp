package testutils

import "testing"

// FailOnError logs the error and terminates the test immediately
func FailOnError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("An error occurred: %v", err)
		t.FailNow()
	}
}