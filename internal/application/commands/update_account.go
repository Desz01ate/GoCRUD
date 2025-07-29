package commands

import (
	"arise_tech_assessment/internal/domain"

	"github.com/google/uuid"
)

type UpdateAccountCommand struct {
	ID         uuid.UUID `json:"id" binding:"required"`
	HolderName string    `json:"holder_name"`
}

type UpdateAccountResponse struct {
	Account *domain.Account `json:"account"`
}
