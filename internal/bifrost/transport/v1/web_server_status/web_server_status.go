package web_server_status

import (
	"github.com/tremendouscan/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/tremendouscan/bifrost/internal/bifrost/transport/v1/options"
)

type webServerStatusServer struct {
	handler handler.WebServerStatusHandlers
	options *options.Options
}

func NewWebServerStatusServer(
	handler handler.WebServerStatusHandlers,
	options *options.Options,
) *webServerStatusServer {
	return &webServerStatusServer{
		handler: handler,
		options: options,
	}
}
