package endpoint

import (
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"

	"github.com/tremendouscan/bifrost/internal/pkg/auth/service"
)

var (
	ErrInvalidLoginReqType  = errors.New("RequestType has only one type: Login")
	ErrInvalidVerifyReqType = errors.New("RequestType has only one type: Verify")
	ErrInvalidLoginRequest  = errors.New("request has only one class: AuthRequest")
	ErrInvalidVerifyRequest = errors.New("request has only one class: VerifyRequest")
)

type AuthRequest struct {
	RequestType string `json:"request_type"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Unexpired   bool   `json:"unexpired"`
}

type AuthResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

type VerifyRequest struct {
	ResquesType string `json:"resques_type"`
	Token       string `json:"token"`
}

type VerifyResponse struct {
	Result bool  `json:"result"`
	Error  error `json:"error"`
}

type AuthEndpoints struct {
	LoginEndpoint       endpoint.Endpoint
	VerifyEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func NewAuthEndpoints(loginEP, verifyEP, healthCheckEP endpoint.Endpoint) AuthEndpoints {
	return AuthEndpoints{
		LoginEndpoint:       loginEP,
		VerifyEndpoint:      verifyEP,
		HealthCheckEndpoint: healthCheckEP,
	}
}

func MakeAuthEndpoints(svc service.Service) AuthEndpoints {
	return NewAuthEndpoints(MakeLoginEndpoint(svc), MakeVerifyEndpoint(svc), MakeHealthCheckEndpoint(svc))
}

func MakeLoginEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(AuthRequest); ok {
			var res string
			if strings.EqualFold(req.RequestType, "Login") {
				res, err = svc.Login(ctx, req.Username, req.Password, req.Unexpired)
			} else {
				return nil, ErrInvalidLoginReqType
			}
			return AuthResponse{
				Result: res,
				Error:  err,
			}, nil
		} else {
			return nil, ErrInvalidLoginRequest
		}
	}
}

func MakeVerifyEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(VerifyRequest); ok {
			var res bool
			if strings.EqualFold(req.ResquesType, "Verify") {
				res, _ = svc.Verify(ctx, req.Token)
			} else {
				return nil, ErrInvalidVerifyReqType
			}
			return VerifyResponse{
				Result: res,
				Error:  err,
			}, nil
		} else {
			return nil, ErrInvalidVerifyRequest
		}
	}
}

// HealthRequest 健康检查请求结构.
type HealthRequest struct{}

// HealthResponse 健康检查响应结构.
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint.
func MakeHealthCheckEndpoint(_ service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthResponse{true}, nil
	}
}
