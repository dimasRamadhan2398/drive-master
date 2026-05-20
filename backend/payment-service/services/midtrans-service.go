// services/midtrans_service.go
package services

import (
	"payment-service/pkg/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService struct {
    snapClient    snap.Client
    coreAPIClient coreapi.Client
    cfg           *config.MidtransConfig
}

func NewMidtransService(cfg *config.MidtransConfig) *MidtransService {
    // determine environment
    env := midtrans.Sandbox
    if cfg.Environment == "production" {
        env = midtrans.Production
    }

    // snap client — for payment page/popup
    snapClient := snap.Client{}
    snapClient.New(cfg.ServerKey, env)

    // core API client — for direct charge and notification
    coreAPIClient := coreapi.Client{}
    coreAPIClient.New(cfg.ServerKey, env)

    return &MidtransService{
        snapClient:    snapClient,
        coreAPIClient: coreAPIClient,
        cfg:           cfg,
    }
}