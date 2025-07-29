package commands

import (
	"arise_tech_assessment/internal/domain"

	"github.com/google/uuid"
)

type CreateTransactionCommand struct {
	Type          domain.TransactionType `json:"type" binding:"required"`
	Amount        domain.Money           `json:"amount" binding:"required"`
	FromAccountID *uuid.UUID             `json:"from_account_id,omitempty"`
	ToAccountID   *uuid.UUID             `json:"to_account_id,omitempty"`
	Description   string                 `json:"description" binding:"required"`
}

type CreateTransactionResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
}
