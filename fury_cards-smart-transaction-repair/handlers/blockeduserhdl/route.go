package blockeduserhdl

import (
	"github.com/melisource/fury_go-core/pkg/web"
)

const (
	unlockPath = "/blockedlist/users/{kycIdentificationID}"
)

type BlockedUserHandlerRouter struct {
	blockedUserHdl *BlockedUserHandler
}

func NewRouter(blockedUserHdl *BlockedUserHandler) *BlockedUserHandlerRouter {
	return &BlockedUserHandlerRouter{
		blockedUserHdl: blockedUserHdl,
	}
}

func (s *BlockedUserHandlerRouter) AddRoutesV1(v1 *web.RouteGroup) {
	rp := v1.Group(unlockPath)
	{
		rp.Post("/unblock", s.blockedUserHdl.UnblockUser, web.UncancellableContext())
	}
}
