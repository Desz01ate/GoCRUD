package handlers

import (
	"arise_tech_assessment/internal/application/commands"
	"arise_tech_assessment/internal/domain"
	"arise_tech_assessment/internal/infrastructure/repository"
	"context"
	"errors"
)

type ProcessTransactionHandler struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

func NewProcessTransactionHandler(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
) *ProcessTransactionHandler {
	return &ProcessTransactionHandler{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (h *ProcessTransactionHandler) Handle(
	ctx context.Context,
	command *commands.ProcessTransactionCommand,
) (*commands.ProcessTransactionResponse, error) {
	transaction, err := h.transactionRepo.GetByID(ctx, command.ID)
	if err != nil {
		return nil, err
	}

	// Check if transaction is pending
	if transaction.Status != domain.TransactionStatusPending {
		return nil, errors.New("transaction is not in pending status")
	}

	switch transaction.Type {
	case domain.TransactionTypeDeposit:
		err = h.processDeposit(ctx, transaction)
	case domain.TransactionTypeWithdraw:
		err = h.processWithdraw(ctx, transaction)
	case domain.TransactionTypeTransfer:
		err = h.processTransfer(ctx, transaction)
	default:
		return nil, errors.New("invalid transaction type")
	}

	if err != nil {
		transaction.Fail()
		h.transactionRepo.Update(ctx, transaction)
		return nil, err
	}

	// Mark transaction as completed
	transaction.Complete()
	err = h.transactionRepo.Update(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return &commands.ProcessTransactionResponse{
		Transaction: transaction,
	}, nil
}

func (h *ProcessTransactionHandler) processDeposit(ctx context.Context, transaction *domain.Transaction) error {
	account, err := h.accountRepo.GetByID(ctx, *transaction.ToAccountID)
	if err != nil {
		return err
	}

	err = account.Credit(transaction.Amount)
	if err != nil {
		return err
	}

	return h.accountRepo.Update(ctx, account)
}

func (h *ProcessTransactionHandler) processWithdraw(ctx context.Context, transaction *domain.Transaction) error {
	account, err := h.accountRepo.GetByID(ctx, *transaction.FromAccountID)
	if err != nil {
		return err
	}

	err = account.Debit(transaction.Amount)
	if err != nil {
		return err
	}

	return h.accountRepo.Update(ctx, account)
}

func (h *ProcessTransactionHandler) processTransfer(ctx context.Context, transaction *domain.Transaction) error {
	// Get both accounts
	fromAccount, err := h.accountRepo.GetByID(ctx, *transaction.FromAccountID)
	if err != nil {
		return err
	}

	toAccount, err := h.accountRepo.GetByID(ctx, *transaction.ToAccountID)
	if err != nil {
		return err
	}

	// Debit from source account
	err = fromAccount.Debit(transaction.Amount)
	if err != nil {
		return err
	}

	// Credit to destination account
	err = toAccount.Credit(transaction.Amount)
	if err != nil {
		return err
	}

	// Update both accounts
	err = h.accountRepo.Update(ctx, fromAccount)
	if err != nil {
		return err
	}

	return h.accountRepo.Update(ctx, toAccount)
}
