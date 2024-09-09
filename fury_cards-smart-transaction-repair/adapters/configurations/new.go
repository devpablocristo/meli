package configurations

import (
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/ports"
	"github.com/melisource/fury_cards-go-toolkit/pkg/configclient/v2"
)

type configService struct {
	serviceCache configclient.ConfigServiceCache
}

func (c *configService) ResetTTLDefaultCache(ttl time.Duration) {
	c.serviceCache = configclient.NewConfigServiceCache(ttl)
}

var _ ports.Configurations = (*configService)(nil)

func New(client configclient.ConfigServiceCache) ports.Configurations {
	return &configService{
		client,
	}
}
