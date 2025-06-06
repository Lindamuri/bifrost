package web_server_status

import svcv1 "github.com/tremendouscan/bifrost/internal/bifrost/service/v1"

type webServerStatusEndpoints struct {
	svc svcv1.ServiceFactory
}

func NewWebServerStatusEndpoints(svc svcv1.ServiceFactory) *webServerStatusEndpoints {
	return &webServerStatusEndpoints{svc: svc}
}
