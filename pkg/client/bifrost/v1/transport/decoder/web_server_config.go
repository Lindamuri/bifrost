package decoder

import (
	"context"

	"github.com/marmotedu/errors"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
	pbv1 "github.com/tremendouscan/bifrost/api/protobuf-spec/bifrostpb/v1"
)

type webConfigServer struct{}

func (w webConfigServer) DecodeResponse(ctx context.Context, resp interface{}) (interface{}, error) {
	switch resp := resp.(type) {
	case *pbv1.ServerNames: // decode `GetServerNames` response
		servernames := make(v1.ServerNames, 0)
		for _, serverName := range resp.Names {
			servernames = append(servernames, v1.ServerName{Name: serverName.GetName()})
		}

		return &servernames, nil
	case *pbv1.ServerConfig: // decode `Get` response
		return &v1.WebServerConfig{
			ServerName:           &v1.ServerName{Name: resp.GetServerName()},
			JsonData:             resp.GetJsonData(),
			OriginalFingerprints: resp.GetOriginalFingerprints(),
		}, nil
	case *pbv1.Response: // decode `Update` response
		return &v1.Response{Message: resp.String()}, nil
	default:
		return nil, errors.Errorf("invalid web server config response: %v", resp)
	}
}

var _ Decoder = webConfigServer{}
