package handlers

import (
	"arise_tech_assessment/internal/application/commands"
	"arise_tech_assessment/internal/domain"
	"arise_tech_assessment/internal/infrastructure/repository"
	"context"
	"errors"
)

type CancelTransactionHandler struct {
	transactionRepo repository.TransactionRepository
}

func NewCancelTransactionHandler(transactionRepo repository.TransactionRepository) *CancelTransactionHandler {
	return &CancelTransactionHandler{
		transactionRepo: transactionRepo,
	}
}

func (h *CancelTransactionHandler) Handle(
	ctx context.Context,
	command *commands.CancelTransactionCommand,
) (*commands.CancelTransactionResponse, error) {
	transaction, err := h.transactionRepo.GetByID(ctx, command.ID)
	if err != nil {
		return nil, err
	}

	// Check if transaction can be cancelled (only pending transactions)
	if transaction.Status != domain.TransactionStatusPending {
		return nil, errors.New("only pending transactions can be cancelled")
	}

	transaction.Cancel()

	err = h.transactionRepo.Update(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return &commands.CancelTransactionResponse{
		Transaction: transaction,
	}, nil
}
