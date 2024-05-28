package globus

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const authBaseUrl = "https://auth.globus.org/v2"

// Returns a two-legged (client credental) http client with oauth2 authentication.
// The function can fail if the token acquisition check fails.
func AuthCreateServiceClient(ctx context.Context, clientID string, clientSecret string, scopes []string) (client *http.Client, err error) {
	conf := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     authBaseUrl + "/oauth2/token",
		Scopes:       scopes,
	}

	// token acquisition check
	_, tokenError := conf.Token(ctx)
	if tokenError != nil {
		return nil, fmt.Errorf("error getting token for client: %s", tokenError.Error())
	}

	return conf.Client(ctx), nil
}

// This is a very basic function that returns an oauth2 config
// with the token url hard-coded to the one provided by Globus.
func AuthGenerateOauthClientConfig(ctx context.Context, clientID string, clientSecret string, redirectURL string, scopes []string) (conf oauth2.Config) {
	conf = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: authBaseUrl + "/oauth2/token",
			AuthURL:  "https://auth.globus.org/v2/oauth2/authorize",
		},
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}

	return conf
}
