package queries

import (
	"arise_tech_assessment/internal/domain"

	"github.com/google/uuid"
)

type GetAccountQuery struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type GetAccountResponse struct {
	Account *domain.Account `json:"account"`
}
