package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SwissOpenEM/globus"
	"golang.org/x/oauth2"
)

func login(authCodeGrant bool, clientID string, clientSecret string, redirectURL string, scopes []string) (*http.Client, error) {
	ctx := context.Background()
	var client *http.Client = nil
	if authCodeGrant {
		// 3-legged OAuth2 authentication - client software authenticates as a user

		// on https://app.globus.org/settings/developers/, you must select "Register a thick client or
		// script that will be installed and run by users on their devices" as registration type

		// create config
		conf := globus.AuthGenerateOauthClientConfig(ctx, clientID, clientSecret, redirectURL, scopes)

		// PKCE verifier
		verifier := oauth2.GenerateVerifier()

		// redirect user to consent page to ask for permission and obtain the code
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
		fmt.Printf("Visit the URL for the auth dialog: %v\n\nEnter the received code here: ", url)

		// read-in and exchange code for token
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			return nil, err
		}
		tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if err != nil {
			return nil, err
		}

		// create client
		client = conf.Client(ctx, tok)
	} else {
		// 2-legged OAuth2 authentication - client software authenticates as itself

		// on https://app.globus.org/settings/developers/, you must select "Register a service account or
		// application for automation" as registration type

		// create client
		var err error
		client, err = globus.AuthCreateServiceClient(ctx, clientID, clientSecret, scopes)
		if err != nil {
			return nil, err
		}
		if client == nil {
			return nil, fmt.Errorf("AUTH error: Client is nil")
		}
	}
	return client, nil
}
