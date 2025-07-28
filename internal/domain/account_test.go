package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewAccount_ShouldInitializeAccountWithCorrectDefaults(t *testing.T) {
	// Arrange
	number := "12345678"
	holderName := "John Doe"
	initialBalance := NewMoney(10000, USD)

	// Act
	account := NewAccount(number, holderName, initialBalance)

	// Assert
	if account.ID == uuid.Nil {
		t.Error("Expected account ID to be generated")
	}
	if account.Number != number {
		t.Errorf("Expected number %s, got %s", number, account.Number)
	}
	if account.HolderName != holderName {
		t.Errorf("Expected holder name %s, got %s", holderName, account.HolderName)
	}
	if account.Balance.Amount != initialBalance.Amount {
		t.Errorf("Expected balance %d, got %d", initialBalance.Amount, account.Balance.Amount)
	}
	if account.Status != AccountStatusActive {
		t.Errorf("Expected status %s, got %s", AccountStatusActive, account.Status)
	}
	if account.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if account.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestAccount_Debit_ShouldCorrectlyDebitAccountBalance(t *testing.T) {
	tests := []struct {
		name           string
		initialBalance Money
		status         AccountStatus
		debitAmount    Money
		expectError    bool
		expectedError  string
		expectedBalance Money
	}{
		{
			"successful debit",
			NewMoney(10000, USD),
			AccountStatusActive,
			NewMoney(5000, USD),
			false,
			"",
			NewMoney(5000, USD),
		},
		{
			"debit from inactive account",
			NewMoney(10000, USD),
			AccountStatusInactive,
			NewMoney(5000, USD),
			true,
			"account is not active",
			NewMoney(10000, USD),
		},
		{
			"debit from blocked account",
			NewMoney(10000, USD),
			AccountStatusBlocked,
			NewMoney(5000, USD),
			true,
			"account is not active",
			NewMoney(10000, USD),
		},
		{
			"insufficient funds",
			NewMoney(5000, USD),
			AccountStatusActive,
			NewMoney(10000, USD),
			true,
			"insufficient funds",
			NewMoney(5000, USD),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			account := &Account{
				ID:         uuid.New(),
				Number:     "12345",
				HolderName: "Test User",
				Balance:    tt.initialBalance,
				Status:     tt.status,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			initialUpdatedAt := account.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure UpdatedAt changes

			// Act
			err := account.Debit(tt.debitAmount)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if account.UpdatedAt.Equal(initialUpdatedAt) {
					t.Error("Expected UpdatedAt to be updated")
				}
			}

			if account.Balance.Amount != tt.expectedBalance.Amount {
				t.Errorf("Expected balance %d, got %d", tt.expectedBalance.Amount, account.Balance.Amount)
			}
		})
	}
}

func TestAccount_Credit_ShouldCorrectlyCreditAccountBalance(t *testing.T) {
	tests := []struct {
		name           string
		initialBalance Money
		status         AccountStatus
		creditAmount   Money
		expectError    bool
		expectedError  string
		expectedBalance Money
	}{
		{
			"successful credit",
			NewMoney(5000, USD),
			AccountStatusActive,
			NewMoney(3000, USD),
			false,
			"",
			NewMoney(8000, USD),
		},
		{
			"credit to inactive account",
			NewMoney(5000, USD),
			AccountStatusInactive,
			NewMoney(3000, USD),
			true,
			"account is not active",
			NewMoney(5000, USD),
		},
		{
			"credit to blocked account",
			NewMoney(5000, USD),
			AccountStatusBlocked,
			NewMoney(3000, USD),
			true,
			"account is not active",
			NewMoney(5000, USD),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			account := &Account{
				ID:         uuid.New(),
				Number:     "12345",
				HolderName: "Test User",
				Balance:    tt.initialBalance,
				Status:     tt.status,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			initialUpdatedAt := account.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure UpdatedAt changes

			// Act
			err := account.Credit(tt.creditAmount)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if account.UpdatedAt.Equal(initialUpdatedAt) {
					t.Error("Expected UpdatedAt to be updated")
				}
			}

			if account.Balance.Amount != tt.expectedBalance.Amount {
				t.Errorf("Expected balance %d, got %d", tt.expectedBalance.Amount, account.Balance.Amount)
			}
		})
	}
}

func TestAccount_Block_ShouldSetAccountStatusToBlocked(t *testing.T) {
	// Arrange
	account := &Account{
		ID:         uuid.New(),
		Number:     "12345",
		HolderName: "Test User",
		Balance:    NewMoney(10000, USD),
		Status:     AccountStatusActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	initialUpdatedAt := account.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure UpdatedAt changes

	// Act
	account.Block()

	// Assert
	if account.Status != AccountStatusBlocked {
		t.Errorf("Expected status %s, got %s", AccountStatusBlocked, account.Status)
	}

	if account.UpdatedAt.Equal(initialUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestAccount_Activate_ShouldSetAccountStatusToActive(t *testing.T) {
	// Arrange
	account := &Account{
		ID:         uuid.New(),
		Number:     "12345",
		HolderName: "Test User",
		Balance:    NewMoney(10000, USD),
		Status:     AccountStatusBlocked,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	initialUpdatedAt := account.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure UpdatedAt changes

	// Act
	account.Activate()

	// Assert
	if account.Status != AccountStatusActive {
		t.Errorf("Expected status %s, got %s", AccountStatusActive, account.Status)
	}

	if account.UpdatedAt.Equal(initialUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}