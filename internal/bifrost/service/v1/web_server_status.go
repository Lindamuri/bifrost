package v1

import (
	"context"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
)

type WebServerStatusService interface {
	Get(ctx context.Context) (*v1.Metrics, error)
}
