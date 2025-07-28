package commands

import (
	"arise_tech_assetment/internal/domain"

	"github.com/google/uuid"
)

type CancelTransactionCommand struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type CancelTransactionResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
}