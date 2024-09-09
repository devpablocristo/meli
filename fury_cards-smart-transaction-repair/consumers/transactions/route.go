package transactions

import "github.com/melisource/fury_go-core/pkg/web"

const (
	transactionsPath = "/consumer/transactions"
)

type TransactionsConsumerRouter struct {
	transactionsConsumer *TransactionsConsumer
}

func NewRouter(transactionsConsumer *TransactionsConsumer) *TransactionsConsumerRouter {
	return &TransactionsConsumerRouter{
		transactionsConsumer: transactionsConsumer,
	}
}

func (s *TransactionsConsumerRouter) AddRoutesV1(v1 *web.RouteGroup) {
	rp := v1.Group(transactionsPath)
	{
		rp.Post("/captures/{siteID}", s.transactionsConsumer.ConsumeCaptures, web.UncancellableContext())
		rp.Post("/reversals/{siteID}", s.transactionsConsumer.ConsumeReversals, web.UncancellableContext())
	}
}
