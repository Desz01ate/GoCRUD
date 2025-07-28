package handlers

import (
	"arise_tech_assetment/internal/application/queries"
	"arise_tech_assetment/internal/infrastructure/repository"
	"context"
)

type GetAccountByNumberHandler struct {
	accountRepo repository.AccountRepository
}

func NewGetAccountByNumberHandler(accountRepo repository.AccountRepository) *GetAccountByNumberHandler {
	return &GetAccountByNumberHandler{
		accountRepo: accountRepo,
	}
}

func (h *GetAccountByNumberHandler) Handle(
	ctx context.Context,
	query *queries.GetAccountByNumberQuery,
) (*queries.GetAccountByNumberResponse, error) {
	account, err := h.accountRepo.FindByNumber(ctx, query.Number)
	if err != nil {
		return nil, err
	}

	return &queries.GetAccountByNumberResponse{
		Account: account,
	}, nil
}