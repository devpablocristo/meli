package validationforcepublishhdl

import (
	"github.com/melisource/fury_go-core/pkg/web"
)

const (
	publishPath = "/validation"
)

type PublishHandlerRouter struct {
	publishhdl *ValidationForcePublishHandler
}

func NewRouter(publishhdl *ValidationForcePublishHandler) *PublishHandlerRouter {
	return &PublishHandlerRouter{
		publishhdl: publishhdl,
	}
}

func (p *PublishHandlerRouter) AddRoutesV1(v1 *web.RouteGroup) {
	rp := v1.Group(publishPath)
	{
		rp.Post("/force", p.publishhdl.Publish, web.UncancellableContext())
	}
}
