package handlers

import (
	"arise_tech_assessment/internal/application/queries"
	"arise_tech_assessment/internal/infrastructure/repository"
	"context"
)

type GetAccountTransactionsHandler struct {
	transactionRepo repository.TransactionRepository
}

func NewGetAccountTransactionsHandler(transactionRepo repository.TransactionRepository) *GetAccountTransactionsHandler {
	return &GetAccountTransactionsHandler{
		transactionRepo: transactionRepo,
	}
}

func (h *GetAccountTransactionsHandler) Handle(
	ctx context.Context,
	query *queries.GetAccountTransactionsQuery,
) (*queries.GetAccountTransactionsResponse, error) {
	req := repository.PaginationRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	pagination, err := h.transactionRepo.FindByAccountIDPaginated(ctx, query.AccountID, req)
	if err != nil {
		return nil, err
	}

	return &queries.GetAccountTransactionsResponse{
		Pagination: pagination,
	}, nil
}
