package encoder

import (
	"context"

	"github.com/marmotedu/errors"

	pbv1 "github.com/tremendouscan/bifrost/api/protobuf-spec/bifrostpb/v1"
)

type webServerStatus struct{}

func (w webServerStatus) EncodeRequest(ctx context.Context, req interface{}) (interface{}, error) {
	switch req := req.(type) {
	case *pbv1.Null:
		return req, nil
	default:
		return nil, errors.Errorf("invalid web server status request: %v", req)
	}
}

var _ Encoder = webServerStatus{}
