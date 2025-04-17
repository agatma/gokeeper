package clients

import (
	"gokeeper/internal/client/core/config"

	"github.com/go-resty/resty/v2"
)

type Clients struct {
	AuthClient    *AuthClient
	PrivateClient *PrivateClient
}

func NewClients(cfg *config.Config) *Clients {
	restyClient := resty.New().
		SetBaseURL("https://" + cfg.Addr).
		SetContentLength(true).
		SetRetryCount(cfg.ServerRetries).
		SetTimeout(cfg.ServerTimeout)
	return &Clients{
		AuthClient:    NewAuthClient(restyClient),
		PrivateClient: NewPrivateClient(restyClient),
	}
}
