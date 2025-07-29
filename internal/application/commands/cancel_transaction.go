package commands

import (
	"arise_tech_assessment/internal/domain"

	"github.com/google/uuid"
)

type CancelTransactionCommand struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type CancelTransactionResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
}
