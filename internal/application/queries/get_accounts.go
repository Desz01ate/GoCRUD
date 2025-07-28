package queries

import (
	"arise_tech_assetment/internal/domain"
	"arise_tech_assetment/internal/infrastructure/repository"
)

type GetAccountsQuery struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GetAccountsResponse struct {
	Pagination *repository.PaginationResponse[domain.Account] `json:"pagination"`
}