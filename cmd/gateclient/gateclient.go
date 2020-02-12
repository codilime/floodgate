package gateclient

import (
	"crypto/tls"
	"net/http"

	"cl-gitlab.intra.codilime.com/spinops/floodgate/config"
	gateapi "cl-gitlab.intra.codilime.com/spinops/floodgate/gateapi"
	"golang.org/x/net/context"
)

type GateapiClient struct {
	// Gate API Client
	*gateapi.APIClient

	// request context
	Context context.Context
}

// NewGateapiClient
// TODO: docs
func NewGateapiClient(floodgateConfig *config.Config) *GateapiClient {
	var gateHTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: floodgateConfig.Insecure},
		},
	}

	cfg := gateapi.NewConfiguration()
	cfg.BasePath = floodgateConfig.Endpoint
	cfg.HTTPClient = gateHTTPClient
	client := gateapi.NewAPIClient(cfg)

	auth := context.WithValue(context.Background(), gateapi.ContextBasicAuth, gateapi.BasicAuth{
		UserName: floodgateConfig.Auth.User,
		Password: floodgateConfig.Auth.Password,
	})

	return &GateapiClient{
		APIClient: client,
		Context:   auth,
	}
}
