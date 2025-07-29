package application

import (
	"arise_tech_assessment/internal/application/handlers"
	"arise_tech_assessment/internal/infrastructure/repository"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterHandlers(db *gorm.DB) {
	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Documentation from https://github.com/mehdihadeli/Go-MediatR/blob/main/readme.md#registering-request-handler-to-the-mediatr
	// is a bit outdated.

	// Register Account Command Handlers
	mediatr.RegisterRequestHandler(
		handlers.NewCreateAccountHandler(accountRepo),
	)

	mediatr.RegisterRequestHandler(
		handlers.NewUpdateAccountHandler(accountRepo),
	)

	mediatr.RegisterRequestHandler(
		handlers.NewDeleteAccountHandler(accountRepo),
	)

	// Register Account Query Handlers
	mediatr.RegisterRequestHandler(
		handlers.NewGetAccountHandler(accountRepo),
	)

	mediatr.RegisterRequestHandler(
		handlers.NewGetAccountsHandler(accountRepo),
	)

	mediatr.RegisterRequestHandler(
		handlers.NewGetAccountByNumberHandler(accountRepo),
	)

	// Register Transaction Command Handlers
	mediatr.RegisterRequestHandler(
		handlers.NewCreateTransactionHandler(transactionRepo, accountRepo),
	)

	mediatr.RegisterRequestHandler(
		handlers.NewProcessTransactionHandler(transactionRepo, accountRepo),
	)

	mediatr.RegisterRequestHandler(
		handlers.NewCancelTransactionHandler(transactionRepo),
	)

	// Register Transaction Query Handlers
	mediatr.RegisterRequestHandler(
		handlers.NewGetTransactionHandler(transactionRepo),
	)

	mediatr.RegisterRequestHandler(
		handlers.NewGetTransactionsHandler(transactionRepo),
	)

	mediatr.RegisterRequestHandler(
		handlers.NewGetAccountTransactionsHandler(transactionRepo),
	)
}
