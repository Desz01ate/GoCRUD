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

func TestCreateAccountHandler_Handle_ShouldSuccessfullyCreateAccount(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateAccountHandler(mockRepo)

	command := &commands.CreateAccountCommand{
		Number:         "12345678",
		HolderName:     "John Doe",
		InitialBalance: domain.NewMoney(10000, domain.USD),
	}

	mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	// Act
	ctx := context.Background()
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
	if account.Number != command.Number {
		t.Errorf("Expected number %s, got %s", command.Number, account.Number)
	}
	if account.HolderName != command.HolderName {
		t.Errorf("Expected holder name %s, got %s", command.HolderName, account.HolderName)
	}
	if account.Balance.Amount != command.InitialBalance.Amount {
		t.Errorf("Expected balance %d, got %d", command.InitialBalance.Amount, account.Balance.Amount)
	}
	if account.Status != domain.AccountStatusActive {
		t.Errorf("Expected status %s, got %s", domain.AccountStatusActive, account.Status)
	}
}

func TestCreateAccountHandler_Handle_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateAccountHandler(mockRepo)

	command := &commands.CreateAccountCommand{
		Number:         "12345678",
		HolderName:     "John Doe",
		InitialBalance: domain.NewMoney(10000, domain.USD),
	}

	mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(errors.New("failed to create account"))

	// Act
	ctx := context.Background()
	response, err := handler.Handle(ctx, command)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
	if err.Error() != "failed to create account" {
		t.Errorf("Expected specific error message, got %s", err.Error())
	}
}

func TestCreateAccountHandler_Handle_ShouldApplyDomainValidationsAndSetDefaults(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewCreateAccountHandler(mockRepo)

	command := &commands.CreateAccountCommand{
		Number:         "12345678",
		HolderName:     "Jane Smith",
		InitialBalance: domain.NewMoney(25000, domain.THB),
	}

	mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)

	// Act
	ctx := context.Background()
	response, err := handler.Handle(ctx, command)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	account := response.Account
	if account.Balance.Currency != domain.THB {
		t.Errorf("Expected currency %s, got %s", domain.THB, account.Balance.Currency)
	}
	if account.ID == uuid.Nil {
		t.Error("Expected account ID to be generated")
	}
	if account.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if account.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}
