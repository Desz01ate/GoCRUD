package handlers

import (
	"arise_tech_assetment/internal/application/commands"
	"arise_tech_assetment/internal/domain"
	"arise_tech_assetment/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionHandler_Handle_ShouldSuccessfullyCreateDepositTransaction(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateTransactionHandler(mockTxRepo, mockAccRepo)
	
	toAccountID := uuid.New()
	command := &commands.CreateTransactionCommand{
		Type:        domain.TransactionTypeDeposit,
		Amount:      domain.NewMoney(5000, domain.USD),
		ToAccountID: &toAccountID,
		Description: "Test deposit",
	}
	ctx := context.Background()
	
	mockTxRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.Transaction")).Return(nil)
	
	// Act
	response, err := handler.Handle(ctx, command)
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	if response.Transaction == nil {
		t.Fatal("Expected transaction in response, got nil")
	}
	
	transaction := response.Transaction
	if transaction.Type != domain.TransactionTypeDeposit {
		t.Errorf("Expected type %s, got %s", domain.TransactionTypeDeposit, transaction.Type)
	}
	
	if transaction.Amount.Amount != command.Amount.Amount {
		t.Errorf("Expected amount %d, got %d", command.Amount.Amount, transaction.Amount.Amount)
	}
	
	if transaction.Status != domain.TransactionStatusPending {
		t.Errorf("Expected status %s, got %s", domain.TransactionStatusPending, transaction.Status)
	}
	
	if transaction.ToAccountID == nil {
		t.Error("Expected ToAccountID to be set")
	} else if *transaction.ToAccountID != toAccountID {
		t.Errorf("Expected ToAccountID %s, got %s", toAccountID, *transaction.ToAccountID)
	}
	
	if transaction.FromAccountID != nil {
		t.Error("Expected FromAccountID to be nil for deposit")
	}
}

func TestCreateTransactionHandler_Handle_ShouldSuccessfullyCreateWithdrawTransaction(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateTransactionHandler(mockTxRepo, mockAccRepo)
	
	fromAccountID := uuid.New()
	command := &commands.CreateTransactionCommand{
		Type:          domain.TransactionTypeWithdraw,
		Amount:        domain.NewMoney(3000, domain.USD),
		FromAccountID: &fromAccountID,
		Description:   "Test withdrawal",
	}
	ctx := context.Background()
	
	mockTxRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.Transaction")).Return(nil)
	
	// Act
	response, err := handler.Handle(ctx, command)
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	transaction := response.Transaction
	if transaction.Type != domain.TransactionTypeWithdraw {
		t.Errorf("Expected type %s, got %s", domain.TransactionTypeWithdraw, transaction.Type)
	}
	
	if transaction.FromAccountID == nil {
		t.Error("Expected FromAccountID to be set")
	} else if *transaction.FromAccountID != fromAccountID {
		t.Errorf("Expected FromAccountID %s, got %s", fromAccountID, *transaction.FromAccountID)
	}
	
	if transaction.ToAccountID != nil {
		t.Error("Expected ToAccountID to be nil for withdrawal")
	}
}

func TestCreateTransactionHandler_Handle_ShouldSuccessfullyCreateTransferTransaction(t *testing.T) {
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateTransactionHandler(mockTxRepo, mockAccRepo)
	
	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	mockTxRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.Transaction")).Return(nil)
	
	command := &commands.CreateTransactionCommand{
		Type:          domain.TransactionTypeTransfer,
		Amount:        domain.NewMoney(7500, domain.USD),
		FromAccountID: &fromAccountID,
		ToAccountID:   &toAccountID,
		Description:   "Test transfer",
	}
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, command)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if response == nil {
		t.Fatal("Expected response, got nil")
	}
	
	transaction := response.Transaction
	if transaction.Type != domain.TransactionTypeTransfer {
		t.Errorf("Expected type %s, got %s", domain.TransactionTypeTransfer, transaction.Type)
	}
	
	if transaction.FromAccountID == nil {
		t.Error("Expected FromAccountID to be set")
	} else if *transaction.FromAccountID != fromAccountID {
		t.Errorf("Expected FromAccountID %s, got %s", fromAccountID, *transaction.FromAccountID)
	}
	
	if transaction.ToAccountID == nil {
		t.Error("Expected ToAccountID to be set")
	} else if *transaction.ToAccountID != toAccountID {
		t.Errorf("Expected ToAccountID %s, got %s", toAccountID, *transaction.ToAccountID)
	}
}

func TestCreateTransactionHandler_Handle_ShouldReturnErrorWhenDepositMissingToAccount(t *testing.T) {
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateTransactionHandler(mockTxRepo, mockAccRepo)
	
	command := &commands.CreateTransactionCommand{
		Type:        domain.TransactionTypeDeposit,
		Amount:      domain.NewMoney(5000, domain.USD),
		ToAccountID: nil, // Missing required field
		Description: "Test deposit",
	}
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, command)
	
	if err == nil {
		t.Error("Expected error for missing ToAccountID, got nil")
	}
	
	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
	
	if err.Error() != "to_account_id is required for deposit" {
		t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

func TestCreateTransactionHandler_Handle_ShouldReturnErrorWhenWithdrawMissingFromAccount(t *testing.T) {
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateTransactionHandler(mockTxRepo, mockAccRepo)
	
	command := &commands.CreateTransactionCommand{
		Type:          domain.TransactionTypeWithdraw,
		Amount:        domain.NewMoney(3000, domain.USD),
		FromAccountID: nil, // Missing required field
		Description:   "Test withdrawal",
	}
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, command)
	
	if err == nil {
		t.Error("Expected error for missing FromAccountID, got nil")
	}
	
	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
	
	if err.Error() != "from_account_id is required for withdrawal" {
		t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

func TestCreateTransactionHandler_Handle_ShouldReturnErrorWhenTransferMissingAccounts(t *testing.T) {
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateTransactionHandler(mockTxRepo, mockAccRepo)
	
	tests := []struct {
		name          string
		fromAccountID *uuid.UUID
		toAccountID   *uuid.UUID
	}{
		{"missing both accounts", nil, nil},
		{"missing from account", nil, func() *uuid.UUID { id := uuid.New(); return &id }()},
		{"missing to account", func() *uuid.UUID { id := uuid.New(); return &id }(), nil},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			command := &commands.CreateTransactionCommand{
				Type:          domain.TransactionTypeTransfer,
				Amount:        domain.NewMoney(5000, domain.USD),
				FromAccountID: tt.fromAccountID,
				ToAccountID:   tt.toAccountID,
				Description:   "Test transfer",
			}
			ctx := context.Background()
			
			// Act
			response, err := handler.Handle(ctx, command)
			
			// Assert
			if err == nil {
				t.Error("Expected error for missing account IDs, got nil")
			}
			
			if response != nil {
				t.Error("Expected nil response on error, got response")
			}
			
			if err.Error() != "both from_account_id and to_account_id are required for transfer" {
				t.Errorf("Expected specific error message, got %s", err.Error())
			}
		})
	}
}

func TestCreateTransactionHandler_Handle_ShouldReturnErrorForInvalidTransactionType(t *testing.T) {
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateTransactionHandler(mockTxRepo, mockAccRepo)
	
	command := &commands.CreateTransactionCommand{
		Type:        "invalid_type",
		Amount:      domain.NewMoney(1000, domain.USD),
		Description: "Test transaction",
	}
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, command)
	
	if err == nil {
		t.Error("Expected error for invalid transaction type, got nil")
	}
	
	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
	
	if err.Error() != "invalid transaction type" {
		t.Errorf("Expected 'invalid transaction type' error, got %s", err.Error())
	}
}

func TestCreateTransactionHandler_Handle_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	mockTxRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.Transaction")).Return(errors.New("failed to create transaction"))
	mockAccRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateTransactionHandler(mockTxRepo, mockAccRepo)
	
	toAccountID := uuid.New()
	command := &commands.CreateTransactionCommand{
		Type:        domain.TransactionTypeDeposit,
		Amount:      domain.NewMoney(5000, domain.USD),
		ToAccountID: &toAccountID,
		Description: "Test deposit",
	}
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, command)
	
	if err == nil {
		t.Error("Expected error from repository, got nil")
	}
	
	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
	
	if err.Error() != "failed to create transaction" {
		t.Errorf("Expected 'failed to create transaction' error, got %s", err.Error())
	}
}