package handlers

import (
	"arise_tech_assetment/internal/application/queries"
	"arise_tech_assetment/internal/domain"
	"arise_tech_assetment/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestGetTransactionHandler_Handle_ShouldSuccessfullyRetrieveTransactionByID(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionHandler(mockRepo)

	txID := uuid.New()
	fromAccountID := uuid.New()
	testTransaction := domain.NewWithdrawTransaction(fromAccountID, domain.NewMoney(7500, domain.USD), "Test withdrawal")
	testTransaction.ID = txID
	testTransaction.SetReference("REF-12345")
	mockRepo.EXPECT().GetByID(mock.Anything, txID).Return(testTransaction, nil)

	query := &queries.GetTransactionQuery{
		ID: txID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, query)

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
	if transaction.ID != txID {
		t.Errorf("Expected transaction ID %s, got %s", txID, transaction.ID)
	}

	if transaction.Type != domain.TransactionTypeWithdraw {
		t.Errorf("Expected type %s, got %s", domain.TransactionTypeWithdraw, transaction.Type)
	}

	if transaction.Amount.Amount != 7500 {
		t.Errorf("Expected amount 7500, got %d", transaction.Amount.Amount)
	}

	if transaction.Amount.Currency != domain.USD {
		t.Errorf("Expected currency %s, got %s", domain.USD, transaction.Amount.Currency)
	}

	if transaction.Status != domain.TransactionStatusPending {
		t.Errorf("Expected status %s, got %s", domain.TransactionStatusPending, transaction.Status)
	}

	if transaction.Description != "Test withdrawal" {
		t.Errorf("Expected description 'Test withdrawal', got %s", transaction.Description)
	}

	if transaction.Reference != "REF-12345" {
		t.Errorf("Expected reference 'REF-12345', got %s", transaction.Reference)
	}

	if transaction.FromAccountID == nil {
		t.Error("Expected FromAccountID to be set")
	} else if *transaction.FromAccountID != fromAccountID {
		t.Errorf("Expected FromAccountID %s, got %s", fromAccountID, *transaction.FromAccountID)
	}
}

func TestGetTransactionHandler_Handle_ShouldReturnErrorWhenTransactionNotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionHandler(mockRepo)

	nonExistentID := uuid.New()
	mockRepo.EXPECT().GetByID(mock.Anything, nonExistentID).Return(nil, errors.New("transaction not found"))

	query := &queries.GetTransactionQuery{
		ID: nonExistentID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, query)

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

func TestGetTransactionHandler_Handle_ShouldRetrieveDifferentTransactionTypes(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionHandler(mockRepo)

	fromAccountID := uuid.New()
	toAccountID := uuid.New()

	tests := []struct {
		name        string
		transaction *domain.Transaction
	}{
		{
			"deposit transaction",
			domain.NewDepositTransaction(toAccountID, domain.NewMoney(5000, domain.THB), "Deposit funds"),
		},
		{
			"withdraw transaction",
			domain.NewWithdrawTransaction(fromAccountID, domain.NewMoney(3000, domain.USD), "ATM withdrawal"),
		},
		{
			"transfer transaction",
			domain.NewTransferTransaction(fromAccountID, toAccountID, domain.NewMoney(2500, domain.USD), "Transfer to friend"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			txID := uuid.New()
			tt.transaction.ID = txID
			mockRepo.EXPECT().GetByID(mock.Anything, txID).Return(tt.transaction, nil)

			query := &queries.GetTransactionQuery{
				ID: txID,
			}
			ctx := context.Background()

			// Act
			response, err := handler.Handle(ctx, query)

			// Assert
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			transaction := response.Transaction
			if transaction.Type != tt.transaction.Type {
				t.Errorf("Expected type %s, got %s", tt.transaction.Type, transaction.Type)
			}

			if transaction.Amount.Amount != tt.transaction.Amount.Amount {
				t.Errorf("Expected amount %d, got %d", tt.transaction.Amount.Amount, transaction.Amount.Amount)
			}

			if transaction.Description != tt.transaction.Description {
				t.Errorf("Expected description %s, got %s", tt.transaction.Description, transaction.Description)
			}
		})
	}
}

func TestGetTransactionHandler_Handle_ShouldRetrieveTransactionWithDifferentStates(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionHandler(mockRepo)

	tests := []struct {
		name   string
		status domain.TransactionStatus
		setup  func(*domain.Transaction)
	}{
		{
			"pending transaction",
			domain.TransactionStatusPending,
			func(tx *domain.Transaction) {
				// Already pending by default
			},
		},
		{
			"completed transaction",
			domain.TransactionStatusCompleted,
			func(tx *domain.Transaction) {
				tx.Complete()
			},
		},
		{
			"failed transaction",
			domain.TransactionStatusFailed,
			func(tx *domain.Transaction) {
				tx.Fail()
			},
		},
		{
			"cancelled transaction",
			domain.TransactionStatusCancelled,
			func(tx *domain.Transaction) {
				tx.Cancel()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			txID := uuid.New()
			transaction := domain.NewTransaction(domain.TransactionTypeDeposit, domain.NewMoney(1000, domain.USD), "Test transaction")
			transaction.ID = txID

			tt.setup(transaction)
			mockRepo.EXPECT().GetByID(mock.Anything, txID).Return(transaction, nil)

			query := &queries.GetTransactionQuery{
				ID: txID,
			}
			ctx := context.Background()

			// Act
			response, err := handler.Handle(ctx, query)

			// Assert
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			returnedTx := response.Transaction
			if returnedTx.Status != tt.status {
				t.Errorf("Expected status %s, got %s", tt.status, returnedTx.Status)
			}

			// Check ProcessedAt is set for completed/failed transactions
			if tt.status == domain.TransactionStatusCompleted || tt.status == domain.TransactionStatusFailed {
				if returnedTx.ProcessedAt == nil {
					t.Error("Expected ProcessedAt to be set for processed transaction")
				}
			} else if tt.status == domain.TransactionStatusCancelled || tt.status == domain.TransactionStatusPending {
				if returnedTx.ProcessedAt != nil && tt.status == domain.TransactionStatusCancelled {
					t.Error("Expected ProcessedAt to be nil for cancelled transaction")
				}
			}
		})
	}
}

func TestGetTransactionHandler_Handle_ShouldRetrieveCorrectTransactionWhenMultipleExist(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionHandler(mockRepo)

	// Create transaction ID for test
	tx2ID := uuid.New()

	// Create the transaction we'll actually query for
	tx2 := domain.NewTransaction(domain.TransactionTypeWithdraw, domain.NewMoney(2000, domain.THB), "Second transaction")
	tx2.ID = tx2ID
	mockRepo.EXPECT().GetByID(mock.Anything, tx2ID).Return(tx2, nil)

	// Query for specific transaction (tx2)
	query := &queries.GetTransactionQuery{
		ID: tx2ID,
	}
	ctx := context.Background()

	// Act
	response, err := handler.Handle(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	transaction := response.Transaction
	if transaction.ID != tx2ID {
		t.Errorf("Expected transaction ID %s, got %s", tx2ID, transaction.ID)
	}

	if transaction.Type != domain.TransactionTypeWithdraw {
		t.Errorf("Expected type %s, got %s", domain.TransactionTypeWithdraw, transaction.Type)
	}

	if transaction.Amount.Currency != domain.THB {
		t.Errorf("Expected currency %s, got %s", domain.THB, transaction.Amount.Currency)
	}

	if transaction.Description != "Second transaction" {
		t.Errorf("Expected description 'Second transaction', got %s", transaction.Description)
	}
}