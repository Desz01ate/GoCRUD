package commands

import (
	"arise_tech_assessment/internal/domain"

	"github.com/google/uuid"
)

type ProcessTransactionCommand struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type ProcessTransactionResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
}
