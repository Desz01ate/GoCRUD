package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

type Transaction struct {
	ID                uuid.UUID         `json:"id" gorm:"type:uuid;primary_key"`
	Type              TransactionType   `json:"type"`
	Status            TransactionStatus `json:"status"`
	Amount            Money             `json:"amount" gorm:"embedded"`
	FromAccountID     *uuid.UUID        `json:"from_account_id,omitempty" gorm:"type:uuid;index"`
	ToAccountID       *uuid.UUID        `json:"to_account_id,omitempty" gorm:"type:uuid;index"`
	FromAccount       *Account          `json:"from_account,omitempty" gorm:"foreignKey:FromAccountID;references:ID"`
	ToAccount         *Account          `json:"to_account,omitempty" gorm:"foreignKey:ToAccountID;references:ID"`
	Description       string            `json:"description"`
	Reference         string            `json:"reference" gorm:"uniqueIndex"`
	ProcessedAt       *time.Time        `json:"processed_at,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

func NewTransaction(txType TransactionType, amount Money, description string) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:          uuid.New(),
		Type:        txType,
		Status:      TransactionStatusPending,
		Amount:      amount,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewDepositTransaction(toAccountID uuid.UUID, amount Money, description string) *Transaction {
	tx := NewTransaction(TransactionTypeDeposit, amount, description)
	tx.ToAccountID = &toAccountID
	return tx
}

func NewWithdrawTransaction(fromAccountID uuid.UUID, amount Money, description string) *Transaction {
	tx := NewTransaction(TransactionTypeWithdraw, amount, description)
	tx.FromAccountID = &fromAccountID
	return tx
}

func NewTransferTransaction(fromAccountID, toAccountID uuid.UUID, amount Money, description string) *Transaction {
	tx := NewTransaction(TransactionTypeTransfer, amount, description)
	tx.FromAccountID = &fromAccountID
	tx.ToAccountID = &toAccountID
	return tx
}

func (t *Transaction) Complete() {
	now := time.Now()
	t.Status = TransactionStatusCompleted
	t.ProcessedAt = &now
	t.UpdatedAt = now
}

func (t *Transaction) Fail() {
	now := time.Now()
	t.Status = TransactionStatusFailed
	t.ProcessedAt = &now
	t.UpdatedAt = now
}

func (t *Transaction) Cancel() {
	t.Status = TransactionStatusCancelled
	t.UpdatedAt = time.Now()
}

func (t *Transaction) SetReference(ref string) {
	t.Reference = ref
	t.UpdatedAt = time.Now()
}