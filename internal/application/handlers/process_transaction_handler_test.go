package handlers

import (
	"arise_tech_assessment/internal/application/commands"
	"arise_tech_assessment/internal/domain"
	"arise_tech_assessment/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestProcessTransactionHandler_Handle_ShouldSuccessfullyProcessDepositTransaction(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewProcessTransactionHandler(mockTxRepo, mockAccRepo)

	accountID := uuid.New()
	account := domain.NewAccount("12345", "John Doe", domain.NewMoney(5000, domain.USD))
	account.ID = accountID

	txID := uuid.New()
	transaction := domain.NewDepositTransaction(accountID, domain.NewMoney(2000, domain.USD), "Deposit")
	transaction.ID = txID

	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
	mockAccRepo.EXPECT().GetByID(mock.Anything, accountID).Return(account, nil)
	mockAccRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)
	mockTxRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == txID && tx.Status == domain.TransactionStatusCompleted
	})).Return(nil)

	command := &commands.ProcessTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Transaction.Status != domain.TransactionStatusCompleted {
		t.Errorf("Expected status %s, got %s", domain.TransactionStatusCompleted, response.Transaction.Status)
	}

	if response.Transaction.ProcessedAt == nil {
		t.Error("Expected ProcessedAt to be set")
	}
}

func TestProcessTransactionHandler_Handle_ShouldSuccessfullyProcessWithdrawTransaction(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewProcessTransactionHandler(mockTxRepo, mockAccRepo)

	accountID := uuid.New()
	account := domain.NewAccount("12345", "John Doe", domain.NewMoney(10000, domain.USD))
	account.ID = accountID

	txID := uuid.New()
	transaction := domain.NewWithdrawTransaction(accountID, domain.NewMoney(3000, domain.USD), "Withdraw")
	transaction.ID = txID

	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
	mockAccRepo.EXPECT().GetByID(mock.Anything, accountID).Return(account, nil)
	mockAccRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)
	mockTxRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == txID && tx.Status == domain.TransactionStatusCompleted
	})).Return(nil)

	command := &commands.ProcessTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Transaction.Status != domain.TransactionStatusCompleted {
		t.Errorf("Expected status %s, got %s", domain.TransactionStatusCompleted, response.Transaction.Status)
	}
}

func TestProcessTransactionHandler_Handle_ShouldSuccessfullyProcessTransferTransaction(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewProcessTransactionHandler(mockTxRepo, mockAccRepo)

	fromAccountID := uuid.New()
	fromAccount := domain.NewAccount("12345", "John Doe", domain.NewMoney(10000, domain.USD))
	fromAccount.ID = fromAccountID

	toAccountID := uuid.New()
	toAccount := domain.NewAccount("67890", "Jane Smith", domain.NewMoney(5000, domain.USD))
	toAccount.ID = toAccountID

	txID := uuid.New()
	transaction := domain.NewTransferTransaction(fromAccountID, toAccountID, domain.NewMoney(2000, domain.USD), "Transfer")
	transaction.ID = txID

	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
	mockAccRepo.EXPECT().GetByID(mock.Anything, fromAccountID).Return(fromAccount, nil)
	mockAccRepo.EXPECT().GetByID(mock.Anything, toAccountID).Return(toAccount, nil)
	mockAccRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil).Times(2)
	mockTxRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == txID && tx.Status == domain.TransactionStatusCompleted
	})).Return(nil)

	command := &commands.ProcessTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Transaction.Status != domain.TransactionStatusCompleted {
		t.Errorf("Expected status %s, got %s", domain.TransactionStatusCompleted, response.Transaction.Status)
	}
}

func TestProcessTransactionHandler_Handle_ShouldReturnErrorForInsufficientFunds(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewProcessTransactionHandler(mockTxRepo, mockAccRepo)

	accountID := uuid.New()
	account := domain.NewAccount("12345", "John Doe", domain.NewMoney(1000, domain.USD))
	account.ID = accountID

	txID := uuid.New()
	transaction := domain.NewWithdrawTransaction(accountID, domain.NewMoney(2000, domain.USD), "Withdraw")
	transaction.ID = txID

	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
	mockAccRepo.EXPECT().GetByID(mock.Anything, accountID).Return(account, nil)
	mockTxRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == txID && tx.Status == domain.TransactionStatusFailed
	})).Return(nil)

	command := &commands.ProcessTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err == nil {
		t.Error("Expected error for insufficient funds, got nil")
	}

	if response != nil {
		t.Error("Expected nil response on error, got response")
	}

	if err.Error() != "insufficient funds" {
		t.Errorf("Expected 'insufficient funds' error, got %s", err.Error())
	}
}

func TestProcessTransactionHandler_Handle_ShouldReturnErrorWhenTransactionNotFound(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewProcessTransactionHandler(mockTxRepo, mockAccRepo)

	nonExistentID := uuid.New()
	mockTxRepo.EXPECT().GetByID(mock.Anything, nonExistentID).Return(nil, errors.New("transaction not found"))

	command := &commands.ProcessTransactionCommand{
		ID: nonExistentID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent transaction, got nil")
	}

	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
}

func TestProcessTransactionHandler_Handle_ShouldReturnErrorWhenTransactionIsNotPending(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewProcessTransactionHandler(mockTxRepo, mockAccRepo)

	txID := uuid.New()
	transaction := domain.NewTransaction(domain.TransactionTypeDeposit, domain.NewMoney(2000, domain.USD), "Deposit")
	transaction.ID = txID
	transaction.Complete()

	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)

	command := &commands.ProcessTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err == nil {
		t.Error("Expected error for non-pending transaction, got nil")
	}

	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
}

func TestProcessTransactionHandler_Handle_ShouldReturnErrorForInactiveAccount(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewProcessTransactionHandler(mockTxRepo, mockAccRepo)

	accountID := uuid.New()
	account := domain.NewAccount("12345", "John Doe", domain.NewMoney(5000, domain.USD))
	account.ID = accountID
	account.Block()

	txID := uuid.New()
	transaction := domain.NewDepositTransaction(accountID, domain.NewMoney(2000, domain.USD), "Deposit")
	transaction.ID = txID

	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
	mockAccRepo.EXPECT().GetByID(mock.Anything, accountID).Return(account, nil)
	mockTxRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == txID && tx.Status == domain.TransactionStatusFailed
	})).Return(nil)

	command := &commands.ProcessTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err == nil {
		t.Error("Expected error for inactive account, got nil")
	}

	if response != nil {
		t.Error("Expected nil response on error, got response")
	}

	if err.Error() != "account is not active" {
		t.Errorf("Expected 'account is not active' error, got %s", err.Error())
	}
}
