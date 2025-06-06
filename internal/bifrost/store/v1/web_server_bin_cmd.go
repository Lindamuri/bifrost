package v1

import (
	"context"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
)

type WebServerBinCMDStore interface {
	Exec(ctx context.Context, request *v1.ExecuteRequest) (*v1.ExecuteResponse, error)
}
