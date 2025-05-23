package endpoint

import (
	"github.com/go-kit/kit/endpoint"

	epv1 "github.com/tremendouscan/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/tremendouscan/bifrost/pkg/client/bifrost/v1/transport"
)

type webServerConfigEndpoints struct {
	transport txpclient.WebServerConfigTransport
}

func (w *webServerConfigEndpoints) EndpointGetServerNames() endpoint.Endpoint {
	return w.transport.GetServerNames().Endpoint()
}

func (w *webServerConfigEndpoints) EndpointGet() endpoint.Endpoint {
	return w.transport.Get().Endpoint()
}

func (w *webServerConfigEndpoints) EndpointUpdate() endpoint.Endpoint {
	return w.transport.Update().Endpoint()
}

func newWebServerConfigEndpoints(factory *factory) epv1.WebServerConfigEndpoints {
	return &webServerConfigEndpoints{transport: factory.transport.WebServerConfig()}
}
