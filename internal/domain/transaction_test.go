package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewTransaction_ShouldInitializeTransactionWithCorrectDefaults(t *testing.T) {
	// Arrange
	txType := TransactionTypeDeposit
	amount := NewMoney(10000, USD)
	description := "Test deposit"

	// Act
	tx := NewTransaction(txType, amount, description)

	// Assert
	if tx.ID == uuid.Nil {
		t.Error("Expected transaction ID to be generated")
	}

	if tx.Type != txType {
		t.Errorf("Expected type %s, got %s", txType, tx.Type)
	}

	if tx.Status != TransactionStatusPending {
		t.Errorf("Expected status %s, got %s", TransactionStatusPending, tx.Status)
	}

	if tx.Amount.Amount != amount.Amount {
		t.Errorf("Expected amount %d, got %d", amount.Amount, tx.Amount.Amount)
	}

	if tx.Description != description {
		t.Errorf("Expected description %s, got %s", description, tx.Description)
	}

	if tx.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if tx.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	if tx.ProcessedAt != nil {
		t.Error("Expected ProcessedAt to be nil for new transaction")
	}
}

func TestNewDepositTransaction_ShouldCreateDepositTransaction(t *testing.T) {
	// Arrange
	accountID := uuid.New()
	amount := NewMoney(10000, USD)
	description := "Deposit transaction"

	// Act
	tx := NewDepositTransaction(accountID, amount, description)

	// Assert
	if tx.Type != TransactionTypeDeposit {
		t.Errorf("Expected type %s, got %s", TransactionTypeDeposit, tx.Type)
	}

	if tx.ToAccountID == nil {
		t.Error("Expected ToAccountID to be set")
	} else if *tx.ToAccountID != accountID {
		t.Errorf("Expected ToAccountID %s, got %s", accountID, *tx.ToAccountID)
	}

	if tx.FromAccountID != nil {
		t.Error("Expected FromAccountID to be nil for deposit")
	}
}

func TestNewWithdrawTransaction_ShouldCreateWithdrawTransaction(t *testing.T) {
	// Arrange
	accountID := uuid.New()
	amount := NewMoney(5000, USD)
	description := "Withdraw transaction"

	// Act
	tx := NewWithdrawTransaction(accountID, amount, description)

	// Assert
	if tx.Type != TransactionTypeWithdraw {
		t.Errorf("Expected type %s, got %s", TransactionTypeWithdraw, tx.Type)
	}

	if tx.FromAccountID == nil {
		t.Error("Expected FromAccountID to be set")
	} else if *tx.FromAccountID != accountID {
		t.Errorf("Expected FromAccountID %s, got %s", accountID, *tx.FromAccountID)
	}

	if tx.ToAccountID != nil {
		t.Error("Expected ToAccountID to be nil for withdraw")
	}
}

func TestNewTransferTransaction_ShouldCreateTransferTransaction(t *testing.T) {
	// Arrange
	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	amount := NewMoney(7500, USD)
	description := "Transfer transaction"

	// Act
	tx := NewTransferTransaction(fromAccountID, toAccountID, amount, description)

	// Assert
	if tx.Type != TransactionTypeTransfer {
		t.Errorf("Expected type %s, got %s", TransactionTypeTransfer, tx.Type)
	}

	if tx.FromAccountID == nil {
		t.Error("Expected FromAccountID to be set")
	} else if *tx.FromAccountID != fromAccountID {
		t.Errorf("Expected FromAccountID %s, got %s", fromAccountID, *tx.FromAccountID)
	}

	if tx.ToAccountID == nil {
		t.Error("Expected ToAccountID to be set")
	} else if *tx.ToAccountID != toAccountID {
		t.Errorf("Expected ToAccountID %s, got %s", toAccountID, *tx.ToAccountID)
	}
}

func TestTransaction_Complete_ShouldSetStatusToCompletedAndProcessedAt(t *testing.T) {
	tx := NewTransaction(TransactionTypeDeposit, NewMoney(1000, USD), "Test")
	
	initialUpdatedAt := tx.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time changes
	
	tx.Complete()
	
	if tx.Status != TransactionStatusCompleted {
		t.Errorf("Expected status %s, got %s", TransactionStatusCompleted, tx.Status)
	}
	
	if tx.ProcessedAt == nil {
		t.Error("Expected ProcessedAt to be set")
	}
	
	if tx.UpdatedAt.Equal(initialUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestTransaction_Fail_ShouldSetStatusToFailedAndProcessedAt(t *testing.T) {
	tx := NewTransaction(TransactionTypeWithdraw, NewMoney(1000, USD), "Test")
	
	initialUpdatedAt := tx.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time changes
	
	tx.Fail()
	
	if tx.Status != TransactionStatusFailed {
		t.Errorf("Expected status %s, got %s", TransactionStatusFailed, tx.Status)
	}
	
	if tx.ProcessedAt == nil {
		t.Error("Expected ProcessedAt to be set")
	}
	
	if tx.UpdatedAt.Equal(initialUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestTransaction_Cancel_ShouldSetStatusToCancelled(t *testing.T) {
	tx := NewTransaction(TransactionTypeTransfer, NewMoney(1000, USD), "Test")
	
	initialUpdatedAt := tx.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time changes
	
	tx.Cancel()
	
	if tx.Status != TransactionStatusCancelled {
		t.Errorf("Expected status %s, got %s", TransactionStatusCancelled, tx.Status)
	}
	
	if tx.ProcessedAt != nil {
		t.Error("Expected ProcessedAt to remain nil for cancelled transaction")
	}
	
	if tx.UpdatedAt.Equal(initialUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestTransaction_SetReference_ShouldSetTransactionReference(t *testing.T) {
	tx := NewTransaction(TransactionTypeDeposit, NewMoney(1000, USD), "Test")
	reference := "REF123456"
	
	initialUpdatedAt := tx.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time changes
	
	tx.SetReference(reference)
	
	if tx.Reference != reference {
		t.Errorf("Expected reference %s, got %s", reference, tx.Reference)
	}
	
	if tx.UpdatedAt.Equal(initialUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}