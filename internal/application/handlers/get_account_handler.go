package handlers

import (
	"arise_tech_assessment/internal/application/queries"
	"arise_tech_assessment/internal/infrastructure/repository"
	"context"
)

type GetAccountHandler struct {
	accountRepository repository.AccountRepository
}

func NewGetAccountHandler(accountRepository repository.AccountRepository) *GetAccountHandler {
	return &GetAccountHandler{
		accountRepository: accountRepository,
	}
}

func (h *GetAccountHandler) Handle(ctx context.Context, query *queries.GetAccountQuery) (*queries.GetAccountResponse, error) {
	account, err := h.accountRepository.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return &queries.GetAccountResponse{
		Account: account,
	}, nil
}
