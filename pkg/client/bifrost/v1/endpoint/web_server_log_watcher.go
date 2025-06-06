package endpoint

import (
	"github.com/go-kit/kit/endpoint"

	epv1 "github.com/tremendouscan/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/tremendouscan/bifrost/pkg/client/bifrost/v1/transport"
)

type webServerLogWatcherEndpoints struct {
	transport txpclient.WebServerLogWatcherTransport
}

func (w *webServerLogWatcherEndpoints) EndpointWatch() endpoint.Endpoint {
	return w.transport.Watch().Endpoint()
}

func newWebServerLogWatcherEndpoints(factory *factory) epv1.WebServerLogWatcherEndpoints {
	return &webServerLogWatcherEndpoints{transport: factory.transport.WebServerLogWatcher()}
}
