package infrastructure

import (
	"arise_tech_assessment/internal/domain"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedData() error {
	log.Println("Starting database seeding...")

	var accountCount int64
	if err := s.db.Model(&domain.Account{}).Count(&accountCount).Error; err != nil {
		return fmt.Errorf("failed to count accounts: %w", err)
	}

	if accountCount > 0 {
		log.Printf("Database already contains %d accounts. Skipping seeding.", accountCount)
		return nil
	}

	accounts := s.createSampleAccounts()

	for _, account := range accounts {
		if err := s.db.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create account %s: %w", account.Number, err)
		}
		log.Printf("Created account: %s (%s)", account.Number, account.HolderName)
	}

	transactions := s.createSampleTransactions(accounts)

	for _, transaction := range transactions {
		if err := s.db.Create(&transaction).Error; err != nil {
			return fmt.Errorf("failed to create transaction %s: %w", transaction.Reference, err)
		}
		log.Printf("Created transaction: %s (%s)", transaction.Reference, transaction.Type)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func (s *Seeder) createSampleAccounts() []domain.Account {
	now := time.Now()

	return []domain.Account{
		{
			ID:         uuid.New(),
			Number:     "ACC001",
			HolderName: "John Doe",
			Balance:    domain.NewMoney(100000, domain.THB),
			Status:     domain.AccountStatusActive,
			CreatedAt:  now.Add(-30 * 24 * time.Hour),
			UpdatedAt:  now.Add(-30 * 24 * time.Hour),
		},
		{
			ID:         uuid.New(),
			Number:     "ACC002",
			HolderName: "Jane Smith",
			Balance:    domain.NewMoney(250000, domain.THB),
			Status:     domain.AccountStatusActive,
			CreatedAt:  now.Add(-25 * 24 * time.Hour),
			UpdatedAt:  now.Add(-25 * 24 * time.Hour),
		},
		{
			ID:         uuid.New(),
			Number:     "ACC003",
			HolderName: "Bob Johnson",
			Balance:    domain.NewMoney(50000, domain.USD),
			Status:     domain.AccountStatusActive,
			CreatedAt:  now.Add(-20 * 24 * time.Hour),
			UpdatedAt:  now.Add(-20 * 24 * time.Hour),
		},
		{
			ID:         uuid.New(),
			Number:     "ACC004",
			HolderName: "Alice Brown",
			Balance:    domain.NewMoney(75000, domain.THB),
			Status:     domain.AccountStatusActive,
			CreatedAt:  now.Add(-15 * 24 * time.Hour),
			UpdatedAt:  now.Add(-15 * 24 * time.Hour),
		},
		{
			ID:         uuid.New(),
			Number:     "ACC005",
			HolderName: "Charlie Wilson",
			Balance:    domain.NewMoney(0, domain.THB),
			Status:     domain.AccountStatusInactive,
			CreatedAt:  now.Add(-10 * 24 * time.Hour),
			UpdatedAt:  now.Add(-10 * 24 * time.Hour),
		},
	}
}

func (s *Seeder) createSampleTransactions(accounts []domain.Account) []domain.Transaction {
	now := time.Now()

	var transactions []domain.Transaction

	depositTx := domain.NewDepositTransaction(
		accounts[0].ID,
		domain.NewMoney(50000, domain.THB),
		"Initial deposit",
	)
	depositTx.Reference = "TXN001"
	depositTx.CreatedAt = now.Add(-29 * 24 * time.Hour)
	depositTx.UpdatedAt = now.Add(-29 * 24 * time.Hour)
	depositTx.Complete()
	transactions = append(transactions, *depositTx)

	// Withdrawal transaction for Jane Smith
	withdrawTx := domain.NewWithdrawTransaction(
		accounts[1].ID,
		domain.NewMoney(25000, domain.THB),
		"ATM withdrawal",
	)
	withdrawTx.Reference = "TXN002"
	withdrawTx.CreatedAt = now.Add(-24 * 24 * time.Hour)
	withdrawTx.UpdatedAt = now.Add(-24 * 24 * time.Hour)
	withdrawTx.Complete()
	transactions = append(transactions, *withdrawTx)

	// Transfer transaction from Jane to John
	transferTx := domain.NewTransferTransaction(
		accounts[1].ID,
		accounts[0].ID,
		domain.NewMoney(10000, domain.THB),
		"Payment for services",
	)
	transferTx.Reference = "TXN003"
	transferTx.CreatedAt = now.Add(-20 * 24 * time.Hour)
	transferTx.UpdatedAt = now.Add(-20 * 24 * time.Hour)
	transferTx.Complete()
	transactions = append(transactions, *transferTx)

	// Pending transaction for Bob Johnson
	pendingTx := domain.NewDepositTransaction(
		accounts[2].ID,
		domain.NewMoney(15000, domain.USD),
		"Pending deposit",
	)
	pendingTx.Reference = "TXN004"
	pendingTx.CreatedAt = now.Add(-2 * 24 * time.Hour)
	pendingTx.UpdatedAt = now.Add(-2 * 24 * time.Hour)
	transactions = append(transactions, *pendingTx)

	// Failed transaction for Alice Brown
	failedTx := domain.NewWithdrawTransaction(
		accounts[3].ID,
		domain.NewMoney(100000, domain.THB),
		"Failed withdrawal attempt",
	)
	failedTx.Reference = "TXN005"
	failedTx.CreatedAt = now.Add(-5 * 24 * time.Hour)
	failedTx.UpdatedAt = now.Add(-5 * 24 * time.Hour)
	failedTx.Fail()
	transactions = append(transactions, *failedTx)

	// Cancelled transaction
	cancelledTx := domain.NewTransferTransaction(
		accounts[0].ID,
		accounts[1].ID,
		domain.NewMoney(5000, domain.THB),
		"Cancelled transfer",
	)
	cancelledTx.Reference = "TXN006"
	cancelledTx.CreatedAt = now.Add(-7 * 24 * time.Hour)
	cancelledTx.UpdatedAt = now.Add(-7 * 24 * time.Hour)
	cancelledTx.Cancel()
	transactions = append(transactions, *cancelledTx)

	return transactions
}
