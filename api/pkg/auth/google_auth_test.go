package auth

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testSetupServiceAccountForGoogleCredentials creates a temporary file containing dummy service account JSON
// then set the environment variable GOOGLE_APPLICATION_CREDENTIALS to point to the file.
//
// This is useful for tests that assume Google Cloud Client libraries can automatically find
// the service account credentials in any environment.
//
// At the end of the test, the returned function can be called to perform cleanup.
func testSetupServiceAccountForGoogleCredentials(t *testing.T) (reset func()) {
	serviceAccountKey := []byte(`{
  "type": "service_account",
  "project_id": "foo",
  "private_key_id": "bar",
  "private_key": "baz",
  "client_email": "foo@example.com",
  "client_id": "bar_client_id",
  "auth_uri": "https://oauth2.googleapis.com/auth",
  "token_uri": "https://oauth2.googleapis.com/token"
}`)

	file, err := os.CreateTemp("", "dummy-service-account")
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(file.Name(), serviceAccountKey, 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", file.Name())
	if err != nil {
		t.Fatal(err)
	}

	return func() {
		err := os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
		if err != nil {
			t.Log("Cleanup failed", err)
		}
		err = os.Remove(file.Name())
		if err != nil {
			t.Log("Cleanup failed", err)
		}
	}
}

func TestInitGoogleClient_NoDefaultCredentials(t *testing.T) {
	client, err := InitGoogleClient(context.Background(), "test.audience")
	assert.NoError(t, err)
	assert.Nil(t, client)
}

func TestInitGoogleClient_ServiceAccount(t *testing.T) {
	reset := testSetupServiceAccountForGoogleCredentials(t)
	defer reset()

	client, err := InitGoogleClient(context.Background(), "test.audience")
	assert.NoError(t, err)
	assert.Nil(t, client)
}
