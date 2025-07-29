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

func TestUpdateAccountHandler_Handle_ShouldSuccessfullyUpdateAccountHolderName(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewUpdateAccountHandler(mockRepo)

	accountID := uuid.New()
	existingAccount := domain.NewAccount("12345", "John Doe", domain.NewMoney(10000, domain.USD))
	existingAccount.ID = accountID

	command := &commands.UpdateAccountCommand{
		ID:         accountID,
		HolderName: "John Smith",
	}
	ctx := context.Background()

	mockRepo.EXPECT().GetByID(mock.Anything, accountID).Return(existingAccount, nil)
	mockRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	// Act
	response, err := handler.Handle(ctx, command)

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
	if account.HolderName != command.HolderName {
		t.Errorf("Expected holder name %s, got %s", command.HolderName, account.HolderName)
	}

	if account.ID != accountID {
		t.Errorf("Expected account ID %s, got %s", accountID, account.ID)
	}
}

func TestUpdateAccountHandler_Handle_ShouldNotUpdateHolderNameIfEmpty(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewUpdateAccountHandler(mockRepo)

	accountID := uuid.New()
	originalHolderName := "John Doe"
	existingAccount := domain.NewAccount("12345", originalHolderName, domain.NewMoney(10000, domain.USD))
	existingAccount.ID = accountID

	command := &commands.UpdateAccountCommand{
		ID:         accountID,
		HolderName: "", // Empty holder name should not update
	}
	ctx := context.Background()

	// Mock GetByID to return existing account
	mockRepo.EXPECT().GetByID(mock.Anything, accountID).Return(existingAccount, nil)
	// Mock Update to simulate successful update (even with empty name, handler decides logic)
	mockRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	account := response.Account
	if account.HolderName != originalHolderName {
		t.Errorf("Expected original holder name %s to remain unchanged, got %s", originalHolderName, account.HolderName)
	}
}

func TestUpdateAccountHandler_Handle_ShouldReturnErrorWhenAccountNotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewUpdateAccountHandler(mockRepo)

	nonExistentID := uuid.New()
	command := &commands.UpdateAccountCommand{
		ID:         nonExistentID,
		HolderName: "John Smith",
	}
	ctx := context.Background()

	// Mock GetByID to return account not found error
	mockRepo.EXPECT().GetByID(mock.Anything, nonExistentID).Return(nil, errors.New("account not found"))

	// Act
	response, err := handler.Handle(ctx, command)

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

func TestUpdateAccountHandler_Handle_ShouldReturnErrorWhenUpdateFails(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewUpdateAccountHandler(mockRepo)

	accountID := uuid.New()
	existingAccount := domain.NewAccount("12345", "John Doe", domain.NewMoney(10000, domain.USD))
	existingAccount.ID = accountID

	command := &commands.UpdateAccountCommand{
		ID:         accountID,
		HolderName: "John Smith",
	}
	ctx := context.Background()

	// Mock GetByID to return existing account
	mockRepo.EXPECT().GetByID(mock.Anything, accountID).Return(existingAccount, nil)
	// Mock Update to return error
	mockRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(errors.New("failed to update account"))

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err == nil {
		t.Error("Expected error from repository update, got nil")
	}

	if response != nil {
		t.Error("Expected nil response on error, got response")
	}

	if err.Error() != "failed to update account" {
		t.Errorf("Expected 'failed to update account' error, got %s", err.Error())
	}
}

func TestUpdateAccountHandler_Handle_ShouldPreserveOtherFieldsWhenUpdatingHolderName(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewUpdateAccountHandler(mockRepo)

	accountID := uuid.New()
	originalBalance := domain.NewMoney(15000, domain.THB)
	originalNumber := "98765"
	existingAccount := domain.NewAccount(originalNumber, "Jane Doe", originalBalance)
	existingAccount.ID = accountID
	existingAccount.Block() // Change status

	command := &commands.UpdateAccountCommand{
		ID:         accountID,
		HolderName: "Jane Smith",
	}
	ctx := context.Background()

	// Mock GetByID to return existing account
	mockRepo.EXPECT().GetByID(mock.Anything, accountID).Return(existingAccount, nil)
	// Mock Update to simulate successful update
	mockRepo.EXPECT().Update(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	account := response.Account

	// Check that only holder name was updated
	if account.HolderName != command.HolderName {
		t.Errorf("Expected holder name %s, got %s", command.HolderName, account.HolderName)
	}

	// Check that other fields were preserved
	if account.Number != originalNumber {
		t.Errorf("Expected number %s to be preserved, got %s", originalNumber, account.Number)
	}

	if account.Balance.Amount != originalBalance.Amount {
		t.Errorf("Expected balance %d to be preserved, got %d", originalBalance.Amount, account.Balance.Amount)
	}

	if account.Balance.Currency != originalBalance.Currency {
		t.Errorf("Expected currency %s to be preserved, got %s", originalBalance.Currency, account.Balance.Currency)
	}

	if account.Status != domain.AccountStatusBlocked {
		t.Errorf("Expected status %s to be preserved, got %s", domain.AccountStatusBlocked, account.Status)
	}
}
