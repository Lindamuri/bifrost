package v1

import (
	"context"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
)

type WebServerLogWatcherService interface {
	Watch(ctx context.Context, request *v1.WebServerLogWatchRequest) (*v1.WebServerLog, error)
}
