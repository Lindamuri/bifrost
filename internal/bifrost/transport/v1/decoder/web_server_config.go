package decoder

import (
	"context"

	"github.com/marmotedu/errors"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
	pbv1 "github.com/tremendouscan/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/tremendouscan/bifrost/internal/pkg/code"
)

type webServerConfig struct{}

var _ Decoder = webServerConfig{}

func (d webServerConfig) DecodeRequest(_ context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *pbv1.Null: // decode `GetServerNames` request
		return r, nil
	case *pbv1.ServerName: // decode `Get` request
		return &v1.ServerName{Name: r.GetName()}, nil
	case *pbv1.ServerConfig: // decode `Update` request
		return &v1.WebServerConfig{
			ServerName:           &v1.ServerName{Name: r.GetServerName()},
			JsonData:             r.GetJsonData(),
			OriginalFingerprints: r.GetOriginalFingerprints(),
		}, nil
	default:
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid request: %v", r)
	}
}

func NewWebServerConfigDecoder() Decoder {
	return new(webServerConfig)
}
