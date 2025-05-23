package nginx

import (
	"context"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"time"

	"github.com/marmotedu/component-base/pkg/version"
	"github.com/shirou/gopsutil/host"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
	"github.com/tremendouscan/bifrost/internal/pkg/monitor"
)

const webServerStatusTimeFormatLayout = "2006/01/02 15:04:05"

type webServerStatusStore struct {
	m                  monitor.Monitor
	webServerInfosFunc func() []*v1.WebServerInfo
	os                 string
	bifrostVersion     string
}

func (w *webServerStatusStore) Get(ctx context.Context) (*v1.Metrics, error) {
	sysInfo := w.m.Report()

	return &v1.Metrics{
		OS:             w.os,
		Time:           time.Now().In(time.Local).Format(webServerStatusTimeFormatLayout),
		Cpu:            sysInfo.CpuUsePct,
		Mem:            sysInfo.MemUsePct,
		Disk:           sysInfo.DiskUsePct,
		StatusList:     w.webServerInfosFunc(),
		BifrostVersion: w.bifrostVersion,
	}, nil
}

func newWebServerStatusStore(store *webServerStore) *webServerStatusStore {
	// get os release info
	var os string
	platform, _, release, err := host.PlatformInformation()
	if err != nil {
		logV1.Warnf("Failed to get platform information. %s", err.Error())
		os = "unknown"
	} else {
		os = platform + " " + release
	}

	return &webServerStatusStore{
		m:                  store.monitor,
		webServerInfosFunc: store.configsManger.GetServerInfos,
		os:                 os,
		bifrostVersion:     version.GitVersion,
	}
}
