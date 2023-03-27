package auth

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testSetupDummyGoogleCredentials creates a temporary file containing dummy credentials JSON
// then set the environment variable GOOGLE_APPLICATION_CREDENTIALS to point to the file.
//
// This is useful for tests that assume Google Cloud Client libraries can automatically find
// the service account credentials in any environment.
//
// At the end of the test, the returned function can be called to perform cleanup.
func testSetupDummyGoogleCredentials(t *testing.T, dummyCredentials []byte) (reset func()) {
	file, err := os.CreateTemp("", "dummy-service-account")
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(file.Name(), dummyCredentials, 0644)
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

func TestTestInitGoogleClient(t *testing.T) {
	// Define tests
	tests := map[string]struct {
		dummyCredential string
		err             string
	}{
		"failure | no default credentials found": {
			err: "fsdfas",
		},
		"failure | invalid json file": {
			dummyCredential: `{`,
			err:             "google: error getting credentials using GOOGLE_APPLICATION_CREDENTIALS environment variable: unexpected end of JSON input",
		},
		"failure | json file with invalid credentials": {
			dummyCredential: `{}`,
			err:             "google: error getting credentials using GOOGLE_APPLICATION_CREDENTIALS environment variable: missing 'type' field in credentials",
		},
		"success | service account": {
			dummyCredential: `{
			    "type": "authorized_user",
			    "project_id": "foo",
			    "private_key_id": "bar",
			    "private_key": "baz",
			    "client_email": "foo@example.com",
			    "client_id": "bar_client_id",
			    "auth_uri": "https://oauth2.googleapis.com/auth",
			    "token_uri": "https://oauth2.googleapis.com/token"
		}`,
		},
		"success | user account": {
			dummyCredential: `{
			    "client_id": "dummyclientid.apps.googleusercontent.com",
			    "client_secret": "dummy-secret",
			    "quota_project_id": "gods-production",
			    "refresh_token": "dummy-token",
			    "type": "authorized_user"
			}`,
		},
	}

	// Run tests
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			if data.dummyCredential != "" {
				reset := testSetupDummyGoogleCredentials(t, []byte(data.dummyCredential))
				defer reset()
			}

			client, err := InitGoogleClient(context.Background(), "test.audience")
			if data.err != "" {
				assert.EqualError(t, err, data.err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

//func TestInitGoogleClient_NoDefaultCredentials(t *testing.T) {
//	client, err := InitGoogleClient(context.Background(), "test.audience")
//	assert.Error(t, err)
//	assert.Nil(t, client)
//}
//
//func TestInitGoogleClient_ServiceAccount(t *testing.T) {
//	reset := testSetupDummyGoogleCredentials(
//		t,
//		[]byte(`{
//  "type": "authorized_user",
//  "project_id": "foo",
//  "private_key_id": "bar",
//  "private_key": "baz",
//  "client_email": "foo@example.com",
//  "client_id": "bar_client_id",
//  "auth_uri": "https://oauth2.googleapis.com/auth",
//  "token_uri": "https://oauth2.googleapis.com/token"
//}`),
//	)
//	defer reset()
//
//	client, err := InitGoogleClient(context.Background(), "test.audience")
//	assert.NoError(t, err)
//	assert.NotNil(t, client)
//}
//
//func TestInitGoogleClient_UserAccount(t *testing.T) {
//	reset := testSetupDummyGoogleCredentials(t,
//		[]byte(`{
//  "client_id": "dummyclientid.apps.googleusercontent.com",
//  "client_secret": "dummy-secret",
//  "quota_project_id": "gods-production",
//  "refresh_token": "dummy-token",
//  "type": "authorized_user"
//}`),
//	)
//	defer reset()
//
//	client, err := InitGoogleClient(context.Background(), "test.audience")
//	assert.NoError(t, err)
//	assert.NotNil(t, client)
//}
