package queries

import (
	"arise_tech_assessment/internal/domain"

	"github.com/google/uuid"
)

type GetTransactionQuery struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type GetTransactionResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
}
