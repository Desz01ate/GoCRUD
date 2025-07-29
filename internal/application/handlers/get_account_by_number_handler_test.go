package handlers

import (
	"arise_tech_assessment/internal/application/queries"
	"arise_tech_assessment/internal/domain"
	"arise_tech_assessment/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestGetAccountByNumberHandler_Handle_ShouldSuccessfullyRetrieveAccountByNumber(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountByNumberHandler(mockRepo)

	accountID := uuid.New()
	accountNumber := "ACC-12345"
	testAccount := domain.NewAccount(accountNumber, "Jane Smith", domain.NewMoney(25000, domain.THB))
	testAccount.ID = accountID

	query := &queries.GetAccountByNumberQuery{
		Number: accountNumber,
	}

	mockRepo.EXPECT().FindByNumber(mock.Anything, accountNumber).Return(testAccount, nil)

	// Act
	ctx := context.Background()
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
	if account.Number != accountNumber {
		t.Errorf("Expected account number %s, got %s", accountNumber, account.Number)
	}

	if account.HolderName != "Jane Smith" {
		t.Errorf("Expected holder name 'Jane Smith', got %s", account.HolderName)
	}

	if account.Balance.Amount != 25000 {
		t.Errorf("Expected balance 25000, got %d", account.Balance.Amount)
	}

	if account.Balance.Currency != domain.THB {
		t.Errorf("Expected currency %s, got %s", domain.THB, account.Balance.Currency)
	}

	if account.ID != accountID {
		t.Errorf("Expected account ID %s, got %s", accountID, account.ID)
	}
}

func TestGetAccountByNumberHandler_Handle_ShouldReturnErrorWhenAccountNotFound(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountByNumberHandler(mockRepo)

	query := &queries.GetAccountByNumberQuery{
		Number: "NON-EXISTENT",
	}

	mockRepo.EXPECT().FindByNumber(mock.Anything, "NON-EXISTENT").Return(nil, errors.New("account not found"))

	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

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

func TestGetAccountByNumberHandler_Handle_ShouldReturnCorrectAccountWhenMultipleExist(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountByNumberHandler(mockRepo)

	targetNumber := "ACC-002"
	targetAccount := domain.NewAccount(targetNumber, "Bob Wilson", domain.NewMoney(10000, domain.USD))
	targetAccount.ID = uuid.New()

	query := &queries.GetAccountByNumberQuery{
		Number: targetNumber,
	}

	mockRepo.EXPECT().FindByNumber(mock.Anything, targetNumber).Return(targetAccount, nil)

	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	account := response.Account
	if account.Number != targetNumber {
		t.Errorf("Expected account number %s, got %s", targetNumber, account.Number)
	}

	if account.HolderName != "Bob Wilson" {
		t.Errorf("Expected holder name 'Bob Wilson', got %s", account.HolderName)
	}
}

func TestGetAccountByNumberHandler_Handle_ShouldReturnErrorForEmptyAccountNumber(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountByNumberHandler(mockRepo)

	query := &queries.GetAccountByNumberQuery{
		Number: "", // Empty number
	}

	mockRepo.EXPECT().FindByNumber(mock.Anything, "").Return(nil, errors.New("account not found"))

	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	if err == nil {
		t.Error("Expected error for empty account number, got nil")
	}

	if response != nil {
		t.Error("Expected nil response on error, got response")
	}

	if err.Error() != "account not found" {
		t.Errorf("Expected 'account not found' error, got %s", err.Error())
	}
}

func TestGetAccountByNumberHandler_Handle_ShouldHandleCaseSensitiveSearch(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountByNumberHandler(mockRepo)

	accountNumber := "ACC-MixedCase123"
	testAccount := domain.NewAccount(accountNumber, "Case Tester", domain.NewMoney(15000, domain.USD))
	testAccount.ID = uuid.New()

	query := &queries.GetAccountByNumberQuery{
		Number: "acc-mixedcase123",
	}

	mockRepo.EXPECT().FindByNumber(mock.Anything, "acc-mixedcase123").Return(nil, errors.New("account not found"))

	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	if err == nil {
		t.Error("Expected error for case-sensitive mismatch, got nil")
	}

	if response != nil {
		t.Error("Expected nil response on case mismatch, got response")
	}

	correctQuery := &queries.GetAccountByNumberQuery{
		Number: accountNumber,
	}

	mockRepo.EXPECT().FindByNumber(mock.Anything, accountNumber).Return(testAccount, nil)

	response, err = handler.Handle(ctx, correctQuery)

	if err != nil {
		t.Errorf("Expected no error with correct case, got %v", err)
	}

	if response == nil || response.Account == nil {
		t.Fatal("Expected account with correct case")
	}

	if response.Account.Number != accountNumber {
		t.Errorf("Expected exact number match %s, got %s", accountNumber, response.Account.Number)
	}
}
