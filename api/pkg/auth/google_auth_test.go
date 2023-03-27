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
	file, err := os.CreateTemp("", "dummy-credentials")
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
			err: "google: could not find default credentials. See " +
				"https://developers.google.com/accounts/docs/application-default-credentials for more information.",
		},
		"failure | invalid json file": {
			dummyCredential: `{`,
			err: "google: error getting credentials using GOOGLE_APPLICATION_CREDENTIALS " +
				"environment variable: unexpected end of JSON input",
		},
		"failure | json file with invalid credentials": {
			dummyCredential: `{}`,
			err: "google: error getting credentials using GOOGLE_APPLICATION_CREDENTIALS " +
				"environment variable: missing 'type' field in credentials",
		},
		"failure | service account not found": {
			//nolint:lll // the private key is in a string literal and cannot contain arbitrary line breaks
			dummyCredential: `{
			    "type": "service_account",
			    "project_id": "example-project",
			    "private_key_id": "1",
			    "private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA4ej0p7bQ7L/r4rVGUz9RN4VQWoej1Bg1mYWIDYslvKrk1gpj\n7wZgkdmM7oVK2OfgrSj/FCTkInKPqaCR0gD7K80q+mLBrN3PUkDrJQZpvRZIff3/\nxmVU1WeruQLFJjnFb2dqu0s/FY/2kWiJtBCakXvXEOb7zfbINuayL+MSsCGSdVYs\nSliS5qQpgyDap+8b5fpXZVJkq92hrcNtbkg7hCYUJczt8n9hcCTJCfUpApvaFQ18\npe+zpyl4+WzkP66I28hniMQyUlA1hBiskT7qiouq0m8IOodhv2fagSZKjOTTU2xk\nSBc//fy3ZpsL7WqgsZS7Q+0VRK8gKfqkxg5OYQIDAQABAoIBAQDGGHzQxGKX+ANk\nnQi53v/c6632dJKYXVJC+PDAz4+bzU800Y+n/bOYsWf/kCp94XcG4Lgsdd0Gx+Zq\nHD9CI1IcqqBRR2AFscsmmX6YzPLTuEKBGMW8twaYy3utlFxElMwoUEsrSWRcCA1y\nnHSDzTt871c7nxCXHxuZ6Nm/XCL7Bg8uidRTSC1sQrQyKgTPhtQdYrPQ4WZ1A4J9\nIisyDYmZodSNZe5P+LTJ6M1SCgH8KH9ZGIxv3diMwzNNpk3kxJc9yCnja4mjiGE2\nYCNusSycU5IhZwVeCTlhQGcNeV/skfg64xkiJE34c2y2ttFbdwBTPixStGaF09nU\nZ422D40BAoGBAPvVyRRsC3BF+qZdaSMFwI1yiXY7vQw5+JZh01tD28NuYdRFzjcJ\nvzT2n8LFpj5ZfZFvSMLMVEFVMgQvWnN0O6xdXvGov6qlRUSGaH9u+TCPNnIldjMP\nB8+xTwFMqI7uQr54wBB+Poq7dVRP+0oHb0NYAwUBXoEuvYo3c/nDoRcZAoGBAOWl\naLHjMv4CJbArzT8sPfic/8waSiLV9Ixs3Re5YREUTtnLq7LoymqB57UXJB3BNz/2\neCueuW71avlWlRtE/wXASj5jx6y5mIrlV4nZbVuyYff0QlcG+fgb6pcJQuO9DxMI\naqFGrWP3zye+LK87a6iR76dS9vRU+bHZpSVvGMKJAoGAFGt3TIKeQtJJyqeUWNSk\nklORNdcOMymYMIlqG+JatXQD1rR6ThgqOt8sgRyJqFCVT++YFMOAqXOBBLnaObZZ\nCFbh1fJ66BlSjoXff0W+SuOx5HuJJAa5+WtFHrPajwxeuRcNa8jwxUsB7n41wADu\nUqWWSRedVBg4Ijbw3nWwYDECgYB0pLew4z4bVuvdt+HgnJA9n0EuYowVdadpTEJg\nsoBjNHV4msLzdNqbjrAqgz6M/n8Ztg8D2PNHMNDNJPVHjJwcR7duSTA6w2p/4k28\nbvvk/45Ta3XmzlxZcZSOct3O31Cw0i2XDVc018IY5be8qendDYM08icNo7vQYkRH\n504kQQKBgQDjx60zpz8ozvm1XAj0wVhi7GwXe+5lTxiLi9Fxq721WDxPMiHDW2XL\nYXfFVy/9/GIMvEiGYdmarK1NW+VhWl1DC5xhDg0kvMfxplt4tynoq1uTsQTY31Mx\nBeF5CT/JuNYk3bEBF0H/Q3VGO1/ggVS+YezdFbLWIRoMnLj6XCFEGg==\n-----END RSA PRIVATE KEY-----\n",
			    "client_email": "service-account@example.com",
			    "client_id": "1234",
			    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
			    "token_uri": "https://accounts.google.com/o/oauth2/token"
			}`,
			err: "oauth2: cannot fetch token: 400 Bad Request\nResponse: " +
				"{\"error\":\"invalid_grant\",\"error_description\":\"Invalid grant: account not found\"}",
		},
		"failure | invalid credentials": {
			dummyCredential: `{
			    "client_id": "dummyclientid.apps.googleusercontent.com",
			    "client_secret": "dummy-secret",
			    "quota_project_id": "gods-production",
			    "refresh_token": "dummy-token",
			    "type": "unauthorized_user"
			}`,
			err: "google: error getting credentials using GOOGLE_APPLICATION_CREDENTIALS environment variable: " +
				"unknown credential type: \"unauthorized_user\"",
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
