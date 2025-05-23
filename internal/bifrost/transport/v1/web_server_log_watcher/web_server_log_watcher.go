package web_server_log_watcher

import (
	"github.com/tremendouscan/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/tremendouscan/bifrost/internal/bifrost/transport/v1/options"
)

type webServerLogWatcherServer struct {
	handler handler.WebServerLogWatcherHandlers
	options *options.Options
}

func NewWebServerLogWatcherServer(
	handler handler.WebServerLogWatcherHandlers,
	opts *options.Options,
) *webServerLogWatcherServer {
	return &webServerLogWatcherServer{
		handler: handler,
		options: opts,
	}
}
