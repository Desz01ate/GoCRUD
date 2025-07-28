package handlers

import (
	"arise_tech_assetment/internal/application/commands"
	"arise_tech_assetment/internal/infrastructure/repository"
	"context"
)

type UpdateAccountHandler struct {
	accountRepo repository.AccountRepository
}

func NewUpdateAccountHandler(accountRepo repository.AccountRepository) *UpdateAccountHandler {
	return &UpdateAccountHandler{
		accountRepo: accountRepo,
	}
}

func (h *UpdateAccountHandler) Handle(
	ctx context.Context,
	command *commands.UpdateAccountCommand,
) (*commands.UpdateAccountResponse, error) {
	account, err := h.accountRepo.GetByID(ctx, command.ID)
	if err != nil {
		return nil, err
	}

	if command.HolderName != "" {
		account.HolderName = command.HolderName
	}

	err = h.accountRepo.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	return &commands.UpdateAccountResponse{
		Account: account,
	}, nil
}
