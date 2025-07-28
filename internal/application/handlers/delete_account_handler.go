package handlers

import (
	"arise_tech_assetment/internal/application/commands"
	"arise_tech_assetment/internal/infrastructure/repository"
	"context"
)

type DeleteAccountHandler struct {
	accountRepo repository.AccountRepository
}

func NewDeleteAccountHandler(accountRepo repository.AccountRepository) *DeleteAccountHandler {
	return &DeleteAccountHandler{
		accountRepo: accountRepo,
	}
}

func (h *DeleteAccountHandler) Handle(
	ctx context.Context,
	command *commands.DeleteAccountCommand,
) (*commands.DeleteAccountResponse, error) {
	err := h.accountRepo.Delete(ctx, command.ID)
	if err != nil {
		return &commands.DeleteAccountResponse{Success: false}, err
	}

	return &commands.DeleteAccountResponse{Success: true}, nil
}