package endpoint

import (
	"github.com/go-kit/kit/endpoint"

	epv1 "github.com/tremendouscan/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/tremendouscan/bifrost/pkg/client/bifrost/v1/transport"
)

type webServerStatisticsEndpoints struct {
	transport txpclient.WebServerStatisticsTransport
}

func (w *webServerStatisticsEndpoints) EndpointGet() endpoint.Endpoint {
	return w.transport.Get().Endpoint()
}

func newWebServerStatisticsEndpoints(factory *factory) epv1.WebServerStatisticsEndpoints {
	return &webServerStatisticsEndpoints{transport: factory.transport.WebServerStatistics()}
}
