package cmd

import (
	"context"
	"fmt"

	"github.com/SwissOpenEM/globus"
	"golang.org/x/oauth2"
)

func login(authCodeGrant bool, clientID string, clientSecret string, redirectURL string, scopes []string) (globus.GlobusClient, error) {
	ctx := context.Background()
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
			return globus.GlobusClient{}, err
		}
		tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if err != nil {
			return globus.GlobusClient{}, err
		}

		// setup auto-refresh & create client
		ts := conf.TokenSource(ctx, tok)
		client := oauth2.NewClient(ctx, ts)

		return globus.HttpClientToGlobusClient(client), nil
	} else {
		// 2-legged OAuth2 authentication - client software authenticates as itself

		// on https://app.globus.org/settings/developers/, you must select "Register a service account or
		// application for automation" as registration type

		// create client
		var err error
		globusClient, err := globus.AuthCreateServiceClient(ctx, clientID, clientSecret, scopes)
		if err != nil {
			return globus.GlobusClient{}, err
		}
		if !globusClient.IsClientSet() {
			return globus.GlobusClient{}, fmt.Errorf("AUTH error: Client is nil")
		}
		return globusClient, nil
	}
}
