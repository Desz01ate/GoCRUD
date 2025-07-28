package commands

import (
	"github.com/google/uuid"
)

type DeleteAccountCommand struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type DeleteAccountResponse struct {
	Success bool `json:"success"`
}