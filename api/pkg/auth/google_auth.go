package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	htransport "google.golang.org/api/transport/http"
)

// JSON key file types.
const (
	serviceAccountKey = "service_account"
)

type credentialsFile struct {
	Type string `json:"type"`
}

// idTokenSource is an oauth2.TokenSource that wraps another TokenSource
// It takes the id_token from TokenSource and passes that on as a bearer token
type idTokenSource struct {
	TokenSource oauth2.TokenSource
}

func (s *idTokenSource) Token() (*oauth2.Token, error) {
	token, err := s.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("token did not contain an id_token")
	}

	return &oauth2.Token{
		AccessToken: idToken,
		TokenType:   "Bearer",
		Expiry:      token.Expiry,
	}, nil
}

// InitGoogleClient is a helper method to be used by CaraML components to initialise a Google Client that appends ID
// tokens to the headers of all outgoing requests with ID tokens, regardless of the type of credentials used
func InitGoogleClient(ctx context.Context, audience string) (*http.Client, error) {
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}

	if len(cred.JSON) < 0 {
		return nil, fmt.Errorf("no default credentials found")
	}

	var f credentialsFile
	if err := json.Unmarshal(cred.JSON, &f); err != nil {
		return nil, err
	}

	if f.Type == serviceAccountKey {
		return idtoken.NewClient(ctx, audience)
	} else {
		defaultTokenSource, err := google.DefaultTokenSource(ctx)
		tokenSource := oauth2.ReuseTokenSource(nil, &idTokenSource{TokenSource: defaultTokenSource})

		var opts []idtoken.ClientOption
		opts = append(opts, option.WithTokenSource(tokenSource), internaloption.SkipDialSettingsValidation())
		t, err := htransport.NewTransport(ctx, http.DefaultTransport, opts...)
		if err != nil {
			return nil, err
		}
		return &http.Client{Transport: t}, nil
	}
}
