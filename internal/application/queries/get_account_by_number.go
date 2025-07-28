package queries

import (
	"arise_tech_assetment/internal/domain"
)

type GetAccountByNumberQuery struct {
	Number string `json:"number" binding:"required"`
}

type GetAccountByNumberResponse struct {
	Account *domain.Account `json:"account"`
}