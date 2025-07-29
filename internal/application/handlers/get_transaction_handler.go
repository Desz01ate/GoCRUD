package handlers

import (
	"arise_tech_assessment/internal/application/queries"
	"arise_tech_assessment/internal/infrastructure/repository"
	"context"
)

type GetTransactionHandler struct {
	transactionRepo repository.TransactionRepository
}

func NewGetTransactionHandler(transactionRepo repository.TransactionRepository) *GetTransactionHandler {
	return &GetTransactionHandler{
		transactionRepo: transactionRepo,
	}
}

func (h *GetTransactionHandler) Handle(
	ctx context.Context,
	query *queries.GetTransactionQuery,
) (*queries.GetTransactionResponse, error) {
	transaction, err := h.transactionRepo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return &queries.GetTransactionResponse{
		Transaction: transaction,
	}, nil
}
