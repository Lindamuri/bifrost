package decoder

import (
	"context"
	"encoding/json"

	"github.com/marmotedu/errors"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
	pbv1 "github.com/tremendouscan/bifrost/api/protobuf-spec/bifrostpb/v1"
)

type webServerStatus struct{}

func (w webServerStatus) DecodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	switch resp := resp.(type) {
	case *pbv1.Metrics:
		metrics := new(v1.Metrics)
		err := json.Unmarshal(resp.GetJsonData(), metrics)

		return metrics, err
	default:
		return nil, errors.Errorf("invalid web server status response: %v", resp)
	}
}

var _ Decoder = webServerStatus{}
