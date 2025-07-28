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

func TestGetAccountHandler_Handle_ShouldSuccessfullyRetrieveAccountByID(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountHandler(mockRepo)

	accountID := uuid.New()
	testAccount := domain.NewAccount("12345", "John Doe", domain.NewMoney(10000, domain.USD))
	testAccount.ID = accountID

	query := &queries.GetAccountQuery{
		ID: accountID,
	}
	ctx := context.Background()

	mockRepo.EXPECT().GetByID(mock.Anything, accountID).Return(testAccount, nil)

	// Act
	response, err := handler.Handle(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Account == nil {
		t.Fatal("Expected account in response, got nil")
	}

	account := response.Account
	if account.ID != accountID {
		t.Errorf("Expected account ID %s, got %s", accountID, account.ID)
	}

	if account.Number != "12345" {
		t.Errorf("Expected account number '12345', got %s", account.Number)
	}

	if account.HolderName != "John Doe" {
		t.Errorf("Expected holder name 'John Doe', got %s", account.HolderName)
	}

	if account.Balance.Amount != 10000 {
		t.Errorf("Expected balance 10000, got %d", account.Balance.Amount)
	}

	if account.Status != domain.AccountStatusActive {
		t.Errorf("Expected status %s, got %s", domain.AccountStatusActive, account.Status)
	}
}

func TestGetAccountHandler_Handle_ShouldReturnErrorWhenAccountNotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountHandler(mockRepo)

	nonExistentID := uuid.New()
	query := &queries.GetAccountQuery{
		ID: nonExistentID,
	}
	ctx := context.Background()

	mockRepo.EXPECT().GetByID(mock.Anything, nonExistentID).Return(nil, errors.New("account not found"))

	// Act
	response, err := handler.Handle(ctx, query)

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent account, got nil")
	}

	if response != nil {
		t.Error("Expected nil response on error, got response")
	}

	if err.Error() != "account not found" {
		t.Errorf("Expected 'account not found' error, got %s", err.Error())
	}
}

func TestGetAccountHandler_Handle_ShouldRetrieveCorrectAccountWhenMultipleExist(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountHandler(mockRepo)

	account1ID := uuid.New()
	account1 := domain.NewAccount("11111", "Alice", domain.NewMoney(5000, domain.USD))
	account1.ID = account1ID

	query := &queries.GetAccountQuery{
		ID: account1ID,
	}
	ctx := context.Background()

	mockRepo.EXPECT().GetByID(mock.Anything, account1ID).Return(account1, nil)

	// Act
	response, err := handler.Handle(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	account := response.Account
	if account.ID != account1ID {
		t.Errorf("Expected account1 ID, got %s", account.ID)
	}

	if account.HolderName != "Alice" {
		t.Errorf("Expected Alice, got %s", account.HolderName)
	}

	if account.Balance.Currency != domain.USD {
		t.Errorf("Expected USD currency, got %s", account.Balance.Currency)
	}
}

func TestGetAccountHandler_Handle_ShouldRetrieveAccountWithDifferentStates(t *testing.T) {

	tests := []struct {
		name     string
		status   domain.AccountStatus
		currency domain.Currency
		amount   int64
	}{
		{"active account", domain.AccountStatusActive, domain.USD, 15000},
		{"blocked account", domain.AccountStatusBlocked, domain.THB, 25000},
		{"inactive account", domain.AccountStatusInactive, domain.USD, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := mocks.NewMockAccountRepository(t)
			handler := NewGetAccountHandler(mockRepo)

			accountID := uuid.New()
			testAccount := domain.NewAccount("98765", "Test User", domain.NewMoney(tt.amount, tt.currency))
			testAccount.ID = accountID

			// Set account status
			switch tt.status {
			case domain.AccountStatusBlocked:
				testAccount.Block()
			case domain.AccountStatusInactive:
				testAccount.Status = domain.AccountStatusInactive
			}

			query := &queries.GetAccountQuery{
				ID: accountID,
			}
			ctx := context.Background()

			mockRepo.EXPECT().GetByID(mock.Anything, accountID).Return(testAccount, nil)

			// Act
			response, err := handler.Handle(ctx, query)

			// Assert
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			account := response.Account
			if account.Status != tt.status {
				t.Errorf("Expected status %s, got %s", tt.status, account.Status)
			}

			if account.Balance.Currency != tt.currency {
				t.Errorf("Expected currency %s, got %s", tt.currency, account.Balance.Currency)
			}

			if account.Balance.Amount != tt.amount {
				t.Errorf("Expected amount %d, got %d", tt.amount, account.Balance.Amount)
			}
		})
	}
}
