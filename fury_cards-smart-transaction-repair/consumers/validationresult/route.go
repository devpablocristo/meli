package validationresult

import "github.com/melisource/fury_go-core/pkg/web"

const (
	transactionsPath = "/consumer/validations"
)

type ValidationResultConsumerRouter struct {
	validationResultConsumer *ValidationResultConsumer
}

func NewRouter(validationResultConsumer *ValidationResultConsumer) *ValidationResultConsumerRouter {
	return &ValidationResultConsumerRouter{
		validationResultConsumer: validationResultConsumer,
	}
}

func (s *ValidationResultConsumerRouter) AddRoutesV1(v1 *web.RouteGroup) {
	rp := v1.Group(transactionsPath)
	{
		rp.Post("", s.validationResultConsumer.ConsumeValidationResult, web.UncancellableContext())
	}
}
