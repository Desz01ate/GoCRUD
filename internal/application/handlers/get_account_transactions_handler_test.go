package handlers

import (
	"arise_tech_assessment/internal/application/queries"
	"arise_tech_assessment/internal/domain"
	"arise_tech_assessment/internal/infrastructure/repository"
	"arise_tech_assessment/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestGetAccountTransactionsHandler_Handle_ShouldSuccessfullyRetrieveAccountTransactions(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetAccountTransactionsHandler(mockRepo)

	accountID := uuid.New()
	testTransactions := make([]domain.Transaction, 3)
	for i := 0; i < 3; i++ {
		transaction := domain.NewDepositTransaction(accountID, domain.NewMoney(int64((i+1)*1000), domain.USD), "Account deposit")
		transaction.ID = uuid.New()
		testTransactions[i] = *transaction
	}

	query := &queries.GetAccountTransactionsQuery{
		AccountID: accountID,
		Page:      1,
		PageSize:  3,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       testTransactions,
		Page:       1,
		PageSize:   3,
		Total:      3,
		TotalPages: 1,
	}
	mockRepo.EXPECT().FindByAccountIDPaginated(mock.Anything, accountID, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 1 && req.PageSize == 3
	})).Return(expectedResponse, nil)

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
	if response.Pagination == nil {
		t.Fatal("Expected pagination in response, got nil")
	}

	pagination := response.Pagination
	if pagination.Page != 1 {
		t.Errorf("Expected page 1, got %d", pagination.Page)
	}
	if pagination.PageSize != 3 {
		t.Errorf("Expected page size 3, got %d", pagination.PageSize)
	}
	if pagination.Total != 3 {
		t.Errorf("Expected total 3, got %d", pagination.Total)
	}
	if len(pagination.Data) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(pagination.Data))
	}

	// Verify all transactions belong to the account
	for _, tx := range pagination.Data {
		if tx.ToAccountID == nil || *tx.ToAccountID != accountID {
			t.Errorf("Expected transaction to belong to account %s", accountID)
		}
	}
}

func TestGetAccountTransactionsHandler_Handle_ShouldRetrieveSecondPageOfTransactions(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetAccountTransactionsHandler(mockRepo)

	accountID := uuid.New()
	testTransactions := make([]domain.Transaction, 2)
	for i := 0; i < 2; i++ {
		transaction := domain.NewWithdrawTransaction(accountID, domain.NewMoney(int64((i+1)*500), domain.USD), "Account withdraw")
		transaction.ID = uuid.New()
		testTransactions[i] = *transaction
	}

	query := &queries.GetAccountTransactionsQuery{
		AccountID: accountID,
		Page:      2,
		PageSize:  3,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       testTransactions,
		Page:       2,
		PageSize:   3,
		Total:      5,
		TotalPages: 2,
	}
	mockRepo.EXPECT().FindByAccountIDPaginated(mock.Anything, accountID, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 2 && req.PageSize == 3
	})).Return(expectedResponse, nil)

	// Act
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	pagination := response.Pagination
	if pagination.Page != 2 {
		t.Errorf("Expected page 2, got %d", pagination.Page)
	}
	if pagination.Total != 5 {
		t.Errorf("Expected total 5, got %d", pagination.Total)
	}
	if pagination.TotalPages != 2 {
		t.Errorf("Expected total pages 2, got %d", pagination.TotalPages)
	}
	if len(pagination.Data) != 2 {
		t.Errorf("Expected 2 transactions in second page, got %d", len(pagination.Data))
	}
}

func TestGetAccountTransactionsHandler_Handle_ShouldReturnEmptyResultWhenNoTransactionsFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetAccountTransactionsHandler(mockRepo)

	accountID := uuid.New()
	query := &queries.GetAccountTransactionsQuery{
		AccountID: accountID,
		Page:      1,
		PageSize:  10,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       []domain.Transaction{},
		Page:       1,
		PageSize:   10,
		Total:      0,
		TotalPages: 0,
	}
	mockRepo.EXPECT().FindByAccountIDPaginated(mock.Anything, accountID, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 1 && req.PageSize == 10
	})).Return(expectedResponse, nil)

	// Act
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	pagination := response.Pagination
	if pagination.Total != 0 {
		t.Errorf("Expected total 0, got %d", pagination.Total)
	}
	if pagination.TotalPages != 0 {
		t.Errorf("Expected total pages 0, got %d", pagination.TotalPages)
	}
	if len(pagination.Data) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(pagination.Data))
	}
}

func TestGetAccountTransactionsHandler_Handle_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetAccountTransactionsHandler(mockRepo)

	accountID := uuid.New()
	query := &queries.GetAccountTransactionsQuery{
		AccountID: accountID,
		Page:      1,
		PageSize:  10,
	}

	mockRepo.EXPECT().FindByAccountIDPaginated(mock.Anything, accountID, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 1 && req.PageSize == 10
	})).Return(nil, errors.New("database connection failed"))

	// Act
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	// Assert
	if err == nil {
		t.Error("Expected error from repository, got nil")
	}
	if response != nil {
		t.Error("Expected nil response on error, got response")
	}
	if err.Error() != "database connection failed" {
		t.Errorf("Expected 'database connection failed' error, got %s", err.Error())
	}
}

func TestGetAccountTransactionsHandler_Handle_ShouldApplyDefaultPaginationWhenNotProvided(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetAccountTransactionsHandler(mockRepo)

	accountID := uuid.New()
	testTransactions := make([]domain.Transaction, 5)
	for i := 0; i < 5; i++ {
		fromAccount := uuid.New()
		transaction := domain.NewTransferTransaction(fromAccount, accountID, domain.NewMoney(int64((i+1)*1000), domain.USD), "Transfer to account")
		transaction.ID = uuid.New()
		testTransactions[i] = *transaction
	}

	query := &queries.GetAccountTransactionsQuery{
		AccountID: accountID,
		Page:      0, // Should default to 1
		PageSize:  0, // Should default to 10
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       testTransactions,
		Page:       1,
		PageSize:   10,
		Total:      5,
		TotalPages: 1,
	}
	mockRepo.EXPECT().FindByAccountIDPaginated(mock.Anything, accountID, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 0 && req.PageSize == 0
	})).Return(expectedResponse, nil)

	// Act
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	pagination := response.Pagination
	if pagination.Page != 1 {
		t.Errorf("Expected page to default to 1, got %d", pagination.Page)
	}
	if pagination.PageSize != 10 {
		t.Errorf("Expected page size to default to 10, got %d", pagination.PageSize)
	}
	if len(pagination.Data) != 5 {
		t.Errorf("Expected 5 transactions with default page size, got %d", len(pagination.Data))
	}
}

func TestGetAccountTransactionsHandler_Handle_ShouldRetrieveMixedTransactionTypes(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetAccountTransactionsHandler(mockRepo)

	accountID := uuid.New()
	otherAccountID := uuid.New()

	deposit := domain.NewDepositTransaction(accountID, domain.NewMoney(2000, domain.USD), "Deposit")
	deposit.ID = uuid.New()

	withdraw := domain.NewWithdrawTransaction(accountID, domain.NewMoney(1000, domain.USD), "Withdraw")
	withdraw.ID = uuid.New()

	transfer := domain.NewTransferTransaction(otherAccountID, accountID, domain.NewMoney(1500, domain.USD), "Transfer")
	transfer.ID = uuid.New()

	testTransactions := []domain.Transaction{*deposit, *withdraw, *transfer}

	query := &queries.GetAccountTransactionsQuery{
		AccountID: accountID,
		Page:      1,
		PageSize:  10,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       testTransactions,
		Page:       1,
		PageSize:   10,
		Total:      3,
		TotalPages: 1,
	}
	mockRepo.EXPECT().FindByAccountIDPaginated(mock.Anything, accountID, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 1 && req.PageSize == 10
	})).Return(expectedResponse, nil)

	// Act
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	pagination := response.Pagination
	if len(pagination.Data) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(pagination.Data))
	}

	// Verify we have different transaction types
	types := make(map[domain.TransactionType]bool)
	for _, tx := range pagination.Data {
		types[tx.Type] = true
	}

	expectedTypes := []domain.TransactionType{
		domain.TransactionTypeDeposit,
		domain.TransactionTypeWithdraw,
		domain.TransactionTypeTransfer,
	}

	for _, expectedType := range expectedTypes {
		if !types[expectedType] {
			t.Errorf("Expected to find transaction type %s", expectedType)
		}
	}
}
