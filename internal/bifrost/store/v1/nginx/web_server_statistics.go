package nginx

import (
	"context"

	"github.com/marmotedu/errors"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
	storev1 "github.com/tremendouscan/bifrost/internal/bifrost/store/v1"
	"github.com/tremendouscan/bifrost/internal/pkg/code"
	"github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration"
)

type webServerStatisticsStore struct {
	statisticians map[string]configuration.Statistician
}

func (w *webServerStatisticsStore) Get(ctx context.Context, servername *v1.ServerName) (*v1.Statistics, error) {
	if statistician, has := w.statisticians[servername.Name]; has {
		return statistician.Statistics(), nil
	}

	return nil, errors.WithCode(code.ErrConfigurationNotFound, "nginx server config '%s' not found", servername.Name)
}

var _ storev1.WebServerStatisticsStore = &webServerStatisticsStore{}

func newNginxStatisticsStore(store *webServerStore) storev1.WebServerStatisticsStore {
	statisticians := make(map[string]configuration.Statistician)
	for servername, config := range store.configsManger.GetConfigs() {
		statisticians[servername] = configuration.NewStatistician(config)
	}

	return &webServerStatisticsStore{statisticians: statisticians}
}
