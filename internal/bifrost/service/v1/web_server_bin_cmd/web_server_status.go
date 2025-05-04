package web_server_bin_cmd

import storev1 "github.com/tremendouscan/bifrost/internal/bifrost/store/v1"

type webServerBinCMDService struct {
	store storev1.StoreFactory
}

func NewWebServerBinCMDService(store storev1.StoreFactory) *webServerBinCMDService {
	return &webServerBinCMDService{store: store}
}
