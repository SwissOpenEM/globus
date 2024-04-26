package globus_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	globustransferrequest "github.com/SwissOpenEM/globus-transfer-request"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	globusClient, globusToken := globustransferrequest.AuthCreateServiceClient(
		ctx,
		"client id here",
		"client secret here",
		[]string{
			"openid",
			"email",
			"profile",
			"urn:globus:auth:scope:transfer.api.globus.org:all",
		},
	)
	if globusClient == nil || globusToken == nil {
		t.Errorf("Got nil for globus token or client")
		t.FailNow()
	}
	resp, err := globusClient.Get("https://auth.globus.org/v2/api/identities")
	if err != nil {
		t.Errorf("Got error for test API call: %s", err.Error())
		t.FailNow()
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respBytes := buf.String()
	fmt.Print(string(respBytes))
}
