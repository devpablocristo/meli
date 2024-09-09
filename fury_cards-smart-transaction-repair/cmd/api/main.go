package main

import (
	"log"

	"github.com/melisource/cards-smart-transaction-repair/handlers/validationforcepublishhdl"

	"github.com/melisource/cards-smart-transaction-repair/cmd/di"

	"github.com/melisource/cards-smart-transaction-repair/handlers/blockeduserhdl"
	"github.com/melisource/cards-smart-transaction-repair/handlers/pinghdl"
	"github.com/melisource/cards-smart-transaction-repair/handlers/reversehdl"

	transactionscons "github.com/melisource/cards-smart-transaction-repair/consumers/transactions"
	validationresultcons "github.com/melisource/cards-smart-transaction-repair/consumers/validationresult"
	"github.com/melisource/fury_go-core/pkg/web"
	"github.com/melisource/fury_go-platform/pkg/fury"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app, err := fury.NewWebApplication(fury.WithErrorHandler(web.DefaultErrorHandler))
	if err != nil {
		return err
	}

	router := app.Group("/cards/smart-transaction-repair")
	v1 := router.Group("/v1")

	pinghdl.NewRouter().AddRoutePing(router)

	cmds := di.ConfigReverseDI()

	reversehdl.NewRouter(cmds.ReverseHandler).AddRoutesV1(v1)
	transactionscons.NewRouter(cmds.TransactionsConsumer).AddRoutesV1(v1)
	blockeduserhdl.NewRouter(cmds.BlockedUserHandler).AddRoutesV1(v1)
	validationresultcons.NewRouter(cmds.ValidationResultConsumer).AddRoutesV1(v1)
	validationforcepublishhdl.NewRouter(cmds.ValidationForcePublishHandler).AddRoutesV1(v1)
	return app.Run()
}
