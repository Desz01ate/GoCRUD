package handlers

import (
	"arise_tech_assetment/internal/application/commands"
	"arise_tech_assetment/internal/domain"
	"arise_tech_assetment/internal/infrastructure/repository"
	"context"
	"errors"
)

type CreateTransactionHandler struct {
	transactionRepository repository.TransactionRepository
	accountRepository     repository.AccountRepository
}

func NewCreateTransactionHandler(transactionRepository repository.TransactionRepository, accountRepository repository.AccountRepository) *CreateTransactionHandler {
	return &CreateTransactionHandler{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
	}
}

func (h *CreateTransactionHandler) Handle(ctx context.Context, command *commands.CreateTransactionCommand) (*commands.CreateTransactionResponse, error) {
	var transaction *domain.Transaction

	switch command.Type {
	case domain.TransactionTypeDeposit:
		if command.ToAccountID == nil {
			return nil, errors.New("to_account_id is required for deposit")
		}
		transaction = domain.NewDepositTransaction(*command.ToAccountID, command.Amount, command.Description)

	case domain.TransactionTypeWithdraw:
		if command.FromAccountID == nil {
			return nil, errors.New("from_account_id is required for withdrawal")
		}
		transaction = domain.NewWithdrawTransaction(*command.FromAccountID, command.Amount, command.Description)

	case domain.TransactionTypeTransfer:
		if command.FromAccountID == nil || command.ToAccountID == nil {
			return nil, errors.New("both from_account_id and to_account_id are required for transfer")
		}
		transaction = domain.NewTransferTransaction(*command.FromAccountID, *command.ToAccountID, command.Amount, command.Description)

	default:
		return nil, errors.New("invalid transaction type")
	}

	if err := h.transactionRepository.Create(ctx, transaction); err != nil {
		return nil, err
	}

	return &commands.CreateTransactionResponse{
		Transaction: transaction,
	}, nil
}
