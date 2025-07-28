package queries

import (
	"arise_tech_assetment/internal/domain"
	"arise_tech_assetment/internal/infrastructure/repository"

	"github.com/google/uuid"
)

type GetAccountTransactionsQuery struct {
	AccountID uuid.UUID `json:"account_id" binding:"required"`
	Page      int       `json:"page"`
	PageSize  int       `json:"page_size"`
}

type GetAccountTransactionsResponse struct {
	Pagination *repository.PaginationResponse[domain.Transaction] `json:"pagination"`
}