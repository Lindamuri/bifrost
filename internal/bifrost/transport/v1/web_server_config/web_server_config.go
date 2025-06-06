package web_server_config

import (
	pbv1 "github.com/tremendouscan/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/tremendouscan/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/tremendouscan/bifrost/internal/bifrost/transport/v1/options"
)

var _ pbv1.WebServerConfigServer = &webServerConfigServer{}

type webServerConfigServer struct {
	handler handler.WebServerConfigHandlers
	options *options.Options
}

func NewWebServerConfigServer(
	handler handler.WebServerConfigHandlers,
	options *options.Options,
) pbv1.WebServerConfigServer {
	return &webServerConfigServer{
		handler: handler,
		options: options,
	}
}
