package reversehdl

import (
	"github.com/melisource/fury_go-core/pkg/web"
)

const (
	reversePath = "/reverse"
)

type ReverseHandlerRouter struct {
	reversehdl *ReverseHandler
}

func NewRouter(reversehdl *ReverseHandler) *ReverseHandlerRouter {
	return &ReverseHandlerRouter{
		reversehdl: reversehdl,
	}
}

func (s *ReverseHandlerRouter) AddRoutesV1(v1 *web.RouteGroup) {
	rp := v1.Group(reversePath)
	{
		rp.Post("/{id}", s.reversehdl.Reverse, web.UncancellableContext())
	}
}
