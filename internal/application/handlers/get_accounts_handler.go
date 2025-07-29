package handlers

import (
	"arise_tech_assessment/internal/application/queries"
	"arise_tech_assessment/internal/infrastructure/repository"
	"context"
)

type GetAccountsHandler struct {
	accountRepository repository.AccountRepository
}

func NewGetAccountsHandler(accountRepository repository.AccountRepository) *GetAccountsHandler {
	return &GetAccountsHandler{
		accountRepository: accountRepository,
	}
}

func (h *GetAccountsHandler) Handle(ctx context.Context, query *queries.GetAccountsQuery) (*queries.GetAccountsResponse, error) {
	req := repository.PaginationRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	pagination, err := h.accountRepository.GetPaginated(ctx, req)
	if err != nil {
		return nil, err
	}

	return &queries.GetAccountsResponse{
		Pagination: pagination,
	}, nil
}
