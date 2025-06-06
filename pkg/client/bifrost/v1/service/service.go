package service

import (
	"context"

	epclient "github.com/tremendouscan/bifrost/pkg/client/bifrost/v1/endpoint"
)

var ctxIns context.Context

type Factory interface {
	WebServerConfig() WebServerConfigService
	WebServerStatistics() WebServerStatisticsService
	WebServerStatus() WebServerStatusService
	WebServerLogWatcher() WebServerLogWatcherService
	WebServerBinCMD() WebServerBinCMDService
}

type factory struct {
	eps epclient.Factory
}

func (f *factory) WebServerConfig() WebServerConfigService {
	return newWebServerConfigService(f)
}

func (f *factory) WebServerStatistics() WebServerStatisticsService {
	return newWebServerStatisticsService(f)
}

func (f *factory) WebServerStatus() WebServerStatusService {
	return newWebServerStatusService(f)
}

func (f *factory) WebServerLogWatcher() WebServerLogWatcherService {
	return newWebServerLogWatcherService(f)
}

func (f *factory) WebServerBinCMD() WebServerBinCMDService {
	return newWebServerBinCMDService(f)
}

func New(endpoint epclient.Factory) Factory {
	return &factory{eps: endpoint}
}

func SetContext(ctx context.Context) {
	ctxIns = ctx
}

func GetContext() context.Context {
	if ctxIns == nil {
		return context.Background()
	}

	return ctxIns
}
