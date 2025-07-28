package commands

import (
	"arise_tech_assetment/internal/domain"

	"github.com/google/uuid"
)

type ProcessTransactionCommand struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type ProcessTransactionResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
}