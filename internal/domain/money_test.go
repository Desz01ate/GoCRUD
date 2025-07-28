package domain

import (
	"testing"
)

func TestNewMoney_ShouldCreateMoneyInstanceWithCorrectValues(t *testing.T) {
	// Arrange
	amount := int64(1000)
	currency := USD

	// Act
	money := NewMoney(amount, currency)

	// Assert
	if money.Amount != amount {
		t.Errorf("Expected amount to be %d, got %d", amount, money.Amount)
	}

	if money.Currency != currency {
		t.Errorf("Expected currency to be %s, got %s", currency, money.Currency)
	}
}

func TestMoney_IsZero_ShouldReturnTrueForZeroAmount(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected bool
	}{
		{"zero amount", NewMoney(0, USD), true},
		{"positive amount", NewMoney(100, USD), false},
		{"negative amount", NewMoney(-100, USD), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			money := tt.money

			// Act
			result := money.IsZero()

			// Assert
			if result != tt.expected {
				t.Errorf("IsZero() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMoney_IsNegative_ShouldReturnTrueForNegativeAmount(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected bool
	}{
		{"negative amount", NewMoney(-100, USD), true},
		{"zero amount", NewMoney(0, USD), false},
		{"positive amount", NewMoney(100, USD), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			money := tt.money

			// Act
			got := money.IsNegative()

			// Assert
			if got != tt.expected {
				t.Errorf("IsNegative() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMoney_IsPositive_ShouldReturnTrueForPositiveAmount(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected bool
	}{
		{"positive amount", NewMoney(100, USD), true},
		{"zero amount", NewMoney(0, USD), false},
		{"negative amount", NewMoney(-100, USD), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			money := tt.money

			// Act
			got := money.IsPositive()

			// Assert
			if got != tt.expected {
				t.Errorf("IsPositive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMoney_Add_ShouldCorrectlyAddMoneyOfSameCurrency(t *testing.T) {
	tests := []struct {
		name        string
		money1      Money
		money2      Money
		expected    Money
		expectError bool
	}{
		{
			"same currency addition",
			NewMoney(100, USD),
			NewMoney(200, USD),
			NewMoney(300, USD),
			false,
		},
		{
			"different currency addition",
			NewMoney(100, USD),
			NewMoney(200, THB),
			Money{},
			true,
		},
		{
			"negative addition",
			NewMoney(100, USD),
			NewMoney(-50, USD),
			NewMoney(50, USD),
			false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			money1 := tt.money1
			money2 := tt.money2
			
			// Act
			result, err := money1.Add(money2)
			
			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result.Amount != tt.expected.Amount {
				t.Errorf("Expected amount %d, got %d", tt.expected.Amount, result.Amount)
			}
			
			if result.Currency != tt.expected.Currency {
				t.Errorf("Expected currency %s, got %s", tt.expected.Currency, result.Currency)
			}
		})
	}
}

func TestMoney_Subtract_ShouldCorrectlySubtractMoneyOfSameCurrency(t *testing.T) {
	tests := []struct {
		name        string
		money1      Money
		money2      Money
		expected    Money
		expectError bool
	}{
		{
			"same currency subtraction",
			NewMoney(300, USD),
			NewMoney(100, USD),
			NewMoney(200, USD),
			false,
		},
		{
			"different currency subtraction",
			NewMoney(300, USD),
			NewMoney(100, THB),
			Money{},
			true,
		},
		{
			"result in negative",
			NewMoney(100, USD),
			NewMoney(200, USD),
			NewMoney(-100, USD),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			money1 := tt.money1
			money2 := tt.money2

			// Act
			result, err := money1.Subtract(money2)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Amount != tt.expected.Amount {
				t.Errorf("Expected amount %d, got %d", tt.expected.Amount, result.Amount)
			}

			if result.Currency != tt.expected.Currency {
				t.Errorf("Expected currency %s, got %s", tt.expected.Currency, result.Currency)
			}
		})
	}
}

func TestMoney_String_ShouldReturnFormattedString(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected string
	}{
		{"positive USD", NewMoney(12345, USD), "123.45 USD"},
		{"zero THB", NewMoney(0, THB), "0.00 THB"},
		{"negative USD", NewMoney(-5432, USD), "-54.32 USD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			money := tt.money

			// Act
			got := money.String()

			// Assert
			if got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMoney_ToFloat_ShouldConvertAmountToFloat(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		expected float64
	}{
		{"positive amount", NewMoney(12345, USD), 123.45},
		{"zero amount", NewMoney(0, USD), 0.0},
		{"negative amount", NewMoney(-5432, USD), -54.32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			money := tt.money

			// Act
			got := money.ToFloat()

			// Assert
			if got != tt.expected {
				t.Errorf("ToFloat() = %v, want %v", got, tt.expected)
			}
		})
	}
}