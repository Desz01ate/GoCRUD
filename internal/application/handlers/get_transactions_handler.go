package handlers

import (
	"arise_tech_assetment/internal/application/queries"
	"arise_tech_assetment/internal/infrastructure/repository"
	"context"
)

type GetTransactionsHandler struct {
	transactionRepo repository.TransactionRepository
}

func NewGetTransactionsHandler(transactionRepo repository.TransactionRepository) *GetTransactionsHandler {
	return &GetTransactionsHandler{
		transactionRepo: transactionRepo,
	}
}

func (h *GetTransactionsHandler) Handle(
	ctx context.Context,
	query *queries.GetTransactionsQuery,
) (*queries.GetTransactionsResponse, error) {
	req := repository.PaginationRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	pagination, err := h.transactionRepo.GetPaginated(ctx, req)
	if err != nil {
		return nil, err
	}

	return &queries.GetTransactionsResponse{
		Pagination: pagination,
	}, nil
}