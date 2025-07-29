package handlers

import (
	"arise_tech_assessment/internal/application/commands"
	"arise_tech_assessment/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestDeleteAccountHandler_Handle_ShouldSuccessfullyDeleteAccount(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewDeleteAccountHandler(mockRepo)

	accountID := uuid.New()
	command := &commands.DeleteAccountCommand{
		ID: accountID,
	}
	ctx := context.Background()

	mockRepo.EXPECT().Delete(mock.Anything, accountID).Return(nil)

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}

func TestDeleteAccountHandler_Handle_ShouldReturnErrorWhenAccountNotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewDeleteAccountHandler(mockRepo)

	nonExistentID := uuid.New()
	command := &commands.DeleteAccountCommand{
		ID: nonExistentID,
	}
	ctx := context.Background()

	mockRepo.EXPECT().Delete(mock.Anything, nonExistentID).Return(errors.New("account not found"))

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent account, got nil")
	}

	if response == nil {
		t.Fatal("Expected response even on error, got nil")
	}

	if response.Success {
		t.Error("Expected success to be false on error")
	}

	if err.Error() != "account not found" {
		t.Errorf("Expected 'account not found' error, got %s", err.Error())
	}
}

func TestDeleteAccountHandler_Handle_ShouldReturnErrorWhenDatabaseFails(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewDeleteAccountHandler(mockRepo)

	accountID := uuid.New()
	command := &commands.DeleteAccountCommand{
		ID: accountID,
	}
	ctx := context.Background()

	mockRepo.EXPECT().Delete(mock.Anything, accountID).Return(errors.New("failed to delete account"))

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err == nil {
		t.Error("Expected error from database failure, got nil")
	}

	if response == nil {
		t.Fatal("Expected response even on error, got nil")
	}

	if response.Success {
		t.Error("Expected success to be false on database error")
	}

	if err.Error() != "failed to delete account" {
		t.Errorf("Expected 'failed to delete account' error, got %s", err.Error())
	}
}

func TestDeleteAccountHandler_Handle_ShouldHandleMultipleAccountsCorrectly(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewDeleteAccountHandler(mockRepo)

	account2ID := uuid.New()
	command := &commands.DeleteAccountCommand{
		ID: account2ID,
	}
	ctx := context.Background()

	mockRepo.EXPECT().Delete(mock.Anything, account2ID).Return(nil)

	// Act
	response, err := handler.Handle(ctx, command)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}
}
