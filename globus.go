package globus

import "net/http"

const transferBaseUrl = "https://transfer.api.globusonline.org/v0.10"

type GlobusClient struct {
	client *http.Client
}

func (g GlobusClient) IsClientSet() bool {
	return g.client != nil
}
