package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
	AccountStatusBlocked  AccountStatus = "blocked"
)

type Account struct {
	ID           uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	Number       string        `json:"number" gorm:"uniqueIndex"`
	HolderName   string        `json:"holder_name"`
	Balance      Money         `json:"balance" gorm:"embedded"`
	Status       AccountStatus `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:FromAccountID;references:ID"`
}

func NewAccount(number, holderName string, initialBalance Money) *Account {
	now := time.Now()
	return &Account{
		ID:         uuid.New(),
		Number:     number,
		HolderName: holderName,
		Balance:    initialBalance,
		Status:     AccountStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (a *Account) Debit(amount Money) error {
	if a.Status != AccountStatusActive {
		return errors.New("account is not active")
	}
	
	if a.Balance.Amount < amount.Amount {
		return errors.New("insufficient funds")
	}
	
	a.Balance.Amount -= amount.Amount
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) Credit(amount Money) error {
	if a.Status != AccountStatusActive {
		return errors.New("account is not active")
	}
	
	a.Balance.Amount += amount.Amount
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) Block() {
	a.Status = AccountStatusBlocked
	a.UpdatedAt = time.Now()
}

func (a *Account) Activate() {
	a.Status = AccountStatusActive
	a.UpdatedAt = time.Now()
}