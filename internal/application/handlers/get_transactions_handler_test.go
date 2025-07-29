package handlers

import (
	"arise_tech_assessment/internal/application/queries"
	"arise_tech_assessment/internal/domain"
	"arise_tech_assessment/internal/infrastructure/repository"
	"arise_tech_assessment/mocks"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestGetTransactionsHandler_Handle_ShouldSuccessfullyRetrieveTransactions(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionsHandler(mockRepo)

	testTransactions := make([]domain.Transaction, 3)
	for i := range 3 {
		accountID := uuid.New()
		transaction := domain.NewDepositTransaction(accountID, domain.NewMoney(int64((i+1)*1000), domain.USD), "Test deposit")
		transaction.ID = uuid.New()
		testTransactions[i] = *transaction
	}

	query := &queries.GetTransactionsQuery{
		Page:     1,
		PageSize: 3,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       testTransactions,
		Page:       1,
		PageSize:   3,
		Total:      3,
		TotalPages: 1,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
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
}

func TestGetTransactionsHandler_Handle_ShouldRetrieveSecondPageOfTransactions(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionsHandler(mockRepo)

	testTransactions := make([]domain.Transaction, 2)
	for i := 0; i < 2; i++ {
		accountID := uuid.New()
		transaction := domain.NewWithdrawTransaction(accountID, domain.NewMoney(int64((i+1)*500), domain.USD), "Test withdraw")
		transaction.ID = uuid.New()
		testTransactions[i] = *transaction
	}

	query := &queries.GetTransactionsQuery{
		Page:     2,
		PageSize: 3,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       testTransactions,
		Page:       2,
		PageSize:   3,
		Total:      5,
		TotalPages: 2,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
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

func TestGetTransactionsHandler_Handle_ShouldReturnEmptyResultWhenNoTransactionsFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionsHandler(mockRepo)

	query := &queries.GetTransactionsQuery{
		Page:     1,
		PageSize: 10,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       []domain.Transaction{},
		Page:       1,
		PageSize:   10,
		Total:      0,
		TotalPages: 0,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
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

func TestGetTransactionsHandler_Handle_ShouldApplyDefaultPaginationWhenNotProvided(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionsHandler(mockRepo)

	testTransactions := make([]domain.Transaction, 5)
	for i := 0; i < 5; i++ {
		fromAccountID := uuid.New()
		toAccountID := uuid.New()
		transaction := domain.NewTransferTransaction(fromAccountID, toAccountID, domain.NewMoney(int64((i+1)*1000), domain.USD), "Test transfer")
		transaction.ID = uuid.New()
		testTransactions[i] = *transaction
	}

	query := &queries.GetTransactionsQuery{
		Page:     1,
		PageSize: 10,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       testTransactions,
		Page:       1,
		PageSize:   10,
		Total:      5,
		TotalPages: 1,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
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

func TestGetTransactionsHandler_Handle_ShouldRetrieveSingleTransactionCorrectly(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockTransactionRepository(t)
	handler := NewGetTransactionsHandler(mockRepo)

	accountID := uuid.New()
	transaction := domain.NewDepositTransaction(accountID, domain.NewMoney(50000, domain.THB), "Single transaction")
	transaction.ID = uuid.New()

	query := &queries.GetTransactionsQuery{
		Page:     1,
		PageSize: 5,
	}

	expectedResponse := &repository.PaginationResponse[domain.Transaction]{
		Data:       []domain.Transaction{*transaction},
		Page:       1,
		PageSize:   5,
		Total:      1,
		TotalPages: 1,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 1 && req.PageSize == 5
	})).Return(expectedResponse, nil)

	// Act
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	pagination := response.Pagination
	if pagination.Total != 1 {
		t.Errorf("Expected total 1, got %d", pagination.Total)
	}
	if pagination.TotalPages != 1 {
		t.Errorf("Expected total pages 1, got %d", pagination.TotalPages)
	}
	if len(pagination.Data) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(pagination.Data))
	}

	returnedTransaction := pagination.Data[0]
	if returnedTransaction.Type != domain.TransactionTypeDeposit {
		t.Errorf("Expected transaction type %s, got %s", domain.TransactionTypeDeposit, returnedTransaction.Type)
	}
	if returnedTransaction.Description != "Single transaction" {
		t.Errorf("Expected description 'Single transaction', got %s", returnedTransaction.Description)
	}
}
