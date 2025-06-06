package endpoint

import (
	"github.com/go-kit/kit/endpoint"

	epv1 "github.com/tremendouscan/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/tremendouscan/bifrost/pkg/client/bifrost/v1/transport"
)

type webServerBinCMDEndpoints struct {
	transport txpclient.WebServerBinCMDTransport
}

func (w *webServerBinCMDEndpoints) EndpointExec() endpoint.Endpoint {
	return w.transport.Exec().Endpoint()
}

func newWebServerBinCMDEndpoints(factory *factory) epv1.WebServerBinCMDEndpoints {
	return &webServerBinCMDEndpoints{transport: factory.transport.WebServerBinCMD()}
}
