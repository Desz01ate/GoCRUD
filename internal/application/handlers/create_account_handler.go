package handlers

import (
	"arise_tech_assetment/internal/application/commands"
	"arise_tech_assetment/internal/domain"
	"arise_tech_assetment/internal/infrastructure/repository"
	"context"
)

type CreateAccountHandler struct {
	accountRepository repository.AccountRepository
}

func NewCreateAccountHandler(accountRepository repository.AccountRepository) *CreateAccountHandler {
	return &CreateAccountHandler{
		accountRepository: accountRepository,
	}
}

func (h *CreateAccountHandler) Handle(ctx context.Context, command *commands.CreateAccountCommand) (*commands.CreateAccountResponse, error) {
	account := domain.NewAccount(command.Number, command.HolderName, command.InitialBalance)

	if err := h.accountRepository.Create(ctx, account); err != nil {
		return nil, err
	}

	return &commands.CreateAccountResponse{
		Account: account,
	}, nil
}
