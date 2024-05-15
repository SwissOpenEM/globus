package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/SwissOpenEM/globus-transfer-request"
	"golang.org/x/oauth2"
)

func main() {
	// auth params
	authCodeGrant := flag.Bool("u", false, "enable authorization code based OAuth2 authentication")
	clientID := flag.String("client-id", "", "set client ID of application")
	clientSecret := flag.String("client-secret", "", "set client secret of application")
	authURL := flag.String("auth-url", "", "set auth url (only used in three-legged mode)")

	flag.Parse()

	scopes := []string{
		"urn:globus:auth:scope:transfer.api.globus.org:all",
	}

	ctx := context.Background()
	var client *http.Client = nil
	if *authCodeGrant {
		// 3-legged OAuth2 authentication - client software authenticates as a user

		// on https://app.globus.org/settings/developers/, you must select "Register a thick client or
		// script that will be installed and run by users on their devices" as registration type

		// create config
		conf := globus.AuthGenerateOauthClientConfig(ctx, *clientID, *clientSecret, *authURL, scopes)

		// PKCE verifier
		verifier := oauth2.GenerateVerifier()

		// redirect user to consent page to ask for permission and obtain the code
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
		fmt.Printf("Visit the URL for the auth dialog: %v", url)

		// read-in and exchange code for token
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}
		tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if err != nil {
			log.Fatal(err)
		}

		// create client
		client = conf.Client(ctx, tok)
	} else {
		// 2-legged OAuth2 authentication - client software authenticates as itself

		// on https://app.globus.org/settings/developers/, you must select "Register a service account or
		// application for automation" as registration type

		// create client
		var err error
		client, err = globus.AuthCreateServiceClient(ctx, *clientID, *clientSecret, scopes)
		if err != nil {
			log.Fatal(err)
		}
		if client == nil {
			log.Fatal("AUTH error: Client is nil\n")
		}
	}
	transferList, err := globus.TransferGetTaskList(client, 0, 50)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Result of request: \n")
	for _, transfer := range transferList.Data {
		fmt.Printf("\n%+v\n", transfer)
	}
}
