package v1

import (
	"context"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
)

type WebServerStatisticsService interface {
	Get(ctx context.Context, servername *v1.ServerName) (*v1.Statistics, error)
}
