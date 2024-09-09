package pinghdl

import (
	"net/http"

	"github.com/melisource/cards-smart-transaction-repair/handlers"
	"github.com/melisource/fury_go-core/pkg/web"
)

const pingPath = "/ping"

type PingHandlerRouter struct {
}

func NewRouter() *PingHandlerRouter {
	return &PingHandlerRouter{}
}

func (p *PingHandlerRouter) AddRoutePing(wg *web.RouteGroup) {
	wg.Get(pingPath, func(w http.ResponseWriter, r *http.Request) error {
		return web.EncodeJSON(w, handlers.SimpleMessageResponse{Message: "pong"}, http.StatusOK)
	})
}
