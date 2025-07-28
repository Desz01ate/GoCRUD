package queries

import (
	"arise_tech_assetment/internal/domain"
	"arise_tech_assetment/internal/infrastructure/repository"
)

type GetTransactionsQuery struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GetTransactionsResponse struct {
	Pagination *repository.PaginationResponse[domain.Transaction] `json:"pagination"`
}