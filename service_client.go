package globustransferrequest

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func CreateServiceClient(ctx context.Context, clientID string, clientSecret string, scopes []string) *http.Client {
	conf := clientcredentials.Config{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		TokenURL:       "https://auth.globus.org/v2/oauth2/token",
		Scopes:         scopes,
		EndpointParams: url.Values{},
		AuthStyle:      oauth2.AuthStyleInParams,
	}

	return conf.Client(ctx)
}
