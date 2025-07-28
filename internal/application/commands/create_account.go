package commands

import (
	"arise_tech_assetment/internal/domain"
)

type CreateAccountCommand struct {
	Number         string       `json:"number" binding:"required"`
	HolderName     string       `json:"holder_name" binding:"required"`
	InitialBalance domain.Money `json:"initial_balance" binding:"required"`
}

type CreateAccountResponse struct {
	Account *domain.Account `json:"account"`
}