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

func TestCancelTransactionHandler_Handle_ShouldSuccessfullyCancelPendingTransaction(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	handler := NewCancelTransactionHandler(mockTxRepo)

	txID := uuid.New()
	transaction := domain.NewTransaction(domain.TransactionTypeDeposit, domain.NewMoney(5000, domain.USD), "Test transaction")
	transaction.ID = txID

	command := &commands.CancelTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()

	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
	mockTxRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == txID && tx.Status == domain.TransactionStatusCancelled
	})).Return(nil)

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

	cancelledTx := response.Transaction
	if cancelledTx.Status != domain.TransactionStatusCancelled {
		t.Errorf("Expected status %s, got %s", domain.TransactionStatusCancelled, cancelledTx.Status)
	}

	if cancelledTx.ID != txID {
		t.Errorf("Expected transaction ID %s, got %s", txID, cancelledTx.ID)
	}

	// Verify ProcessedAt is not set for cancelled transactions
	if cancelledTx.ProcessedAt != nil {
		t.Error("Expected ProcessedAt to be nil for cancelled transaction")
	}
}

func TestCancelTransactionHandler_Handle_ShouldReturnErrorWhenTransactionNotFound(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	handler := NewCancelTransactionHandler(mockTxRepo)
	
	nonExistentID := uuid.New()
	command := &commands.CancelTransactionCommand{
		ID: nonExistentID,
	}
	ctx := context.Background()
	
	mockTxRepo.EXPECT().GetByID(mock.Anything, nonExistentID).Return(nil, errors.New("transaction not found"))
	
	// Act
	response, err := handler.Handle(ctx, command)
	
	// Assert
	if err == nil {
		t.Error("Expected error for non-existent transaction, got nil")
	}
	
	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
	
	if err.Error() != "transaction not found" {
		t.Errorf("Expected 'transaction not found' error, got %s", err.Error())
	}
}

func TestCancelTransactionHandler_Handle_ShouldReturnErrorWhenTransactionIsNotPending(t *testing.T) {
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	handler := NewCancelTransactionHandler(mockTxRepo)
	
	tests := []struct {
		name   string
		status domain.TransactionStatus
	}{
		{"completed transaction", domain.TransactionStatusCompleted},
		{"failed transaction", domain.TransactionStatusFailed},
		{"already cancelled transaction", domain.TransactionStatusCancelled},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			txID := uuid.New()
			transaction := domain.NewTransaction(domain.TransactionTypeWithdraw, domain.NewMoney(3000, domain.USD), "Test transaction")
			transaction.ID = txID
			
			// Set transaction to non-pending status
			switch tt.status {
			case domain.TransactionStatusCompleted:
				transaction.Complete()
			case domain.TransactionStatusFailed:
				transaction.Fail()
			case domain.TransactionStatusCancelled:
				transaction.Cancel()
			}
			
			command := &commands.CancelTransactionCommand{
				ID: txID,
			}
			ctx := context.Background()
			
			mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
			
			// Act
			response, err := handler.Handle(ctx, command)
			
			// Assert
			if err == nil {
				t.Error("Expected error for non-pending transaction, got nil")
			}
			
			if response != nil {
				t.Error("Expected nil response on error, got response")
			}
			
			if err.Error() != "only pending transactions can be cancelled" {
				t.Errorf("Expected specific error message, got %s", err.Error())
			}
		})
	}
}

func TestCancelTransactionHandler_Handle_ShouldReturnErrorWhenUpdateFails(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	handler := NewCancelTransactionHandler(mockTxRepo)
	
	txID := uuid.New()
	transaction := domain.NewTransaction(domain.TransactionTypeTransfer, domain.NewMoney(2000, domain.USD), "Test transaction")
	transaction.ID = txID
	
	command := &commands.CancelTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()
	
	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
	mockTxRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == txID && tx.Status == domain.TransactionStatusCancelled
	})).Return(errors.New("failed to update transaction"))
	
	// Act
	response, err := handler.Handle(ctx, command)
	
	// Assert
	if err == nil {
		t.Error("Expected error from repository update, got nil")
	}
	
	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
	
	if err.Error() != "failed to update transaction" {
		t.Errorf("Expected 'failed to update transaction' error, got %s", err.Error())
	}
}

func TestCancelTransactionHandler_Handle_ShouldPreserveTransactionDataUponCancellation(t *testing.T) {
	// Arrange
	mockTxRepo := mocks.NewMockTransactionRepository(t)
	handler := NewCancelTransactionHandler(mockTxRepo)
	
	txID := uuid.New()
	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	originalAmount := domain.NewMoney(7500, domain.THB)
	originalDescription := "Important transfer transaction"
	
	transaction := domain.NewTransferTransaction(fromAccountID, toAccountID, originalAmount, originalDescription)
	transaction.ID = txID
	transaction.SetReference("REF123456")
	
	command := &commands.CancelTransactionCommand{
		ID: txID,
	}
	ctx := context.Background()
	
	mockTxRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)
	mockTxRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == txID && tx.Status == domain.TransactionStatusCancelled
	})).Return(nil)
	
	// Act
	response, err := handler.Handle(ctx, command)
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	cancelledTx := response.Transaction
	
	// Verify only status changed, other fields preserved
	if cancelledTx.Status != domain.TransactionStatusCancelled {
		t.Errorf("Expected status %s, got %s", domain.TransactionStatusCancelled, cancelledTx.Status)
	}
	
	if cancelledTx.Type != domain.TransactionTypeTransfer {
		t.Errorf("Expected type %s to be preserved, got %s", domain.TransactionTypeTransfer, cancelledTx.Type)
	}
	
	if cancelledTx.Amount.Amount != originalAmount.Amount {
		t.Errorf("Expected amount %d to be preserved, got %d", originalAmount.Amount, cancelledTx.Amount.Amount)
	}
	
	if cancelledTx.Amount.Currency != originalAmount.Currency {
		t.Errorf("Expected currency %s to be preserved, got %s", originalAmount.Currency, cancelledTx.Amount.Currency)
	}
	
	if cancelledTx.Description != originalDescription {
		t.Errorf("Expected description %s to be preserved, got %s", originalDescription, cancelledTx.Description)
	}
	
	if cancelledTx.Reference != "REF123456" {
		t.Errorf("Expected reference to be preserved, got %s", cancelledTx.Reference)
	}
	
	if cancelledTx.FromAccountID == nil {
		t.Error("Expected FromAccountID to be preserved")
	} else if *cancelledTx.FromAccountID != fromAccountID {
		t.Errorf("Expected FromAccountID %s to be preserved, got %s", fromAccountID, *cancelledTx.FromAccountID)
	}
	
	if cancelledTx.ToAccountID == nil {
		t.Error("Expected ToAccountID to be preserved")
	} else if *cancelledTx.ToAccountID != toAccountID {
		t.Errorf("Expected ToAccountID %s to be preserved, got %s", toAccountID, *cancelledTx.ToAccountID)
	}
}