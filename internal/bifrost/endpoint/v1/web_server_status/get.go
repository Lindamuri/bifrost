package web_server_status

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/marmotedu/errors"

	pbv1 "github.com/tremendouscan/bifrost/api/protobuf-spec/bifrostpb/v1"
)

func (w *webServerStatusEndpoints) EndpointGet() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if _, ok := request.(*pbv1.Null); ok {
			return w.svc.WebServerStatus().Get(ctx)
		}

		return nil, errors.Errorf("invalid get request, need *pbv1.Null, not %T", request)
	}
}
