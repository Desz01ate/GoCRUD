package handlers

import (
	"arise_tech_assetment/internal/application/queries"
	"arise_tech_assetment/internal/domain"
	"arise_tech_assetment/internal/infrastructure/repository"
	"arise_tech_assetment/mocks"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestGetAccountsHandler_Handle_ShouldSuccessfullyRetrieveAccounts(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountsHandler(mockRepo)

	testAccounts := make([]domain.Account, 5)
	for i := 0; i < 5; i++ {
		account := domain.NewAccount("ACC-"+string(rune('1'+i)), "User "+string(rune('A'+i)), domain.NewMoney(int64((i+1)*1000), domain.USD))
		account.ID = uuid.New()
		testAccounts[i] = *account
	}

	query := &queries.GetAccountsQuery{
		Page:     1,
		PageSize: 3,
	}
	ctx := context.Background()

	expectedResponse := &repository.PaginationResponse[domain.Account]{
		Data:       testAccounts[:3],
		Page:       1,
		PageSize:   3,
		Total:      5,
		TotalPages: 2,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 1 && req.PageSize == 3
	})).Return(expectedResponse, nil)

	// Act
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

	if pagination.Total != 5 {
		t.Errorf("Expected total 5, got %d", pagination.Total)
	}

	if pagination.TotalPages != 2 {
		t.Errorf("Expected total pages 2, got %d", pagination.TotalPages)
	}

	if len(pagination.Data) != 3 {
		t.Errorf("Expected 3 accounts in first page, got %d", len(pagination.Data))
	}
}

func TestGetAccountsHandler_Handle_ShouldRetrieveSecondPageOfAccounts(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountsHandler(mockRepo)
	
	// Create 7 test accounts
	testAccounts := make([]domain.Account, 7)
	for i := 0; i < 7; i++ {
		account := domain.NewAccount("ACC-"+string(rune('1'+i)), "User "+string(rune('A'+i)), domain.NewMoney(int64((i+1)*1000), domain.USD))
		account.ID = uuid.New()
		testAccounts[i] = *account
	}
	
	query := &queries.GetAccountsQuery{
		Page:     2,
		PageSize: 3,
	}
	
	// Mock the GetPaginated call for second page
	expectedResponse := &repository.PaginationResponse[domain.Account]{
		Data:       testAccounts[3:6],
		Page:       2,
		PageSize:   3,
		Total:      7,
		TotalPages: 3,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 2 && req.PageSize == 3
	})).Return(expectedResponse, nil)
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	pagination := response.Pagination
	if pagination.Page != 2 {
		t.Errorf("Expected page 2, got %d", pagination.Page)
	}
	
	if pagination.Total != 7 {
		t.Errorf("Expected total 7, got %d", pagination.Total)
	}
	
	if pagination.TotalPages != 3 {
		t.Errorf("Expected total pages 3, got %d", pagination.TotalPages)
	}
	
	if len(pagination.Data) != 3 {
		t.Errorf("Expected 3 accounts in second page, got %d", len(pagination.Data))
	}
}

func TestGetAccountsHandler_Handle_ShouldRetrieveLastPageOfAccounts(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountsHandler(mockRepo)
	
	// Create 8 test accounts
	testAccounts := make([]domain.Account, 8)
	for i := 0; i < 8; i++ {
		account := domain.NewAccount("ACC-"+string(rune('1'+i)), "User "+string(rune('A'+i)), domain.NewMoney(int64((i+1)*1000), domain.USD))
		account.ID = uuid.New()
		testAccounts[i] = *account
	}
	
	query := &queries.GetAccountsQuery{
		Page:     3,
		PageSize: 3,
	}
	
	// Mock the GetPaginated call for last page
	expectedResponse := &repository.PaginationResponse[domain.Account]{
		Data:       testAccounts[6:8], // Last 2 accounts
		Page:       3,
		PageSize:   3,
		Total:      8,
		TotalPages: 3,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 3 && req.PageSize == 3
	})).Return(expectedResponse, nil)
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	pagination := response.Pagination
	if pagination.Page != 3 {
		t.Errorf("Expected page 3, got %d", pagination.Page)
	}
	
	if pagination.Total != 8 {
		t.Errorf("Expected total 8, got %d", pagination.Total)
	}
	
	if pagination.TotalPages != 3 {
		t.Errorf("Expected total pages 3, got %d", pagination.TotalPages)
	}
	
	// Last page should have only 2 accounts (8 % 3 = 2)
	if len(pagination.Data) != 2 {
		t.Errorf("Expected 2 accounts in last page, got %d", len(pagination.Data))
	}
}

func TestGetAccountsHandler_Handle_ShouldReturnEmptyResultWhenNoAccountsFound(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountsHandler(mockRepo)
	
	query := &queries.GetAccountsQuery{
		Page:     1,
		PageSize: 10,
	}
	
	// Mock empty result
	expectedResponse := &repository.PaginationResponse[domain.Account]{
		Data:       []domain.Account{},
		Page:       1,
		PageSize:   10,
		Total:      0,
		TotalPages: 0,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 1 && req.PageSize == 10
	})).Return(expectedResponse, nil)
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)
	
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
		t.Errorf("Expected 0 accounts, got %d", len(pagination.Data))
	}
}

func TestGetAccountsHandler_Handle_ShouldApplyDefaultPaginationWhenNotProvided(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountsHandler(mockRepo)
	
	// Create 5 test accounts
	testAccounts := make([]domain.Account, 5)
	for i := 0; i < 5; i++ {
		account := domain.NewAccount("ACC-"+string(rune('1'+i)), "User "+string(rune('A'+i)), domain.NewMoney(int64((i+1)*1000), domain.USD))
		account.ID = uuid.New()
		testAccounts[i] = *account
	}
	
	// Test with zero/negative values that should default
	query := &queries.GetAccountsQuery{
		Page:     0, // Should default to 1
		PageSize: 0, // Should default to 10
	}
	
	// Mock the GetPaginated call with defaults
	expectedResponse := &repository.PaginationResponse[domain.Account]{
		Data:       testAccounts,
		Page:       1,
		PageSize:   10,
		Total:      5,
		TotalPages: 1,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 0 && req.PageSize == 0
	})).Return(expectedResponse, nil)
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)
	
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
	
	// Should return all 5 accounts since page size is 10
	if len(pagination.Data) != 5 {
		t.Errorf("Expected 5 accounts with default page size, got %d", len(pagination.Data))
	}
}

func TestGetAccountsHandler_Handle_ShouldReturnEmptyResultWhenPageBeyondRange(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountsHandler(mockRepo)
	
	// Request page 5 when there are only 3 accounts
	query := &queries.GetAccountsQuery{
		Page:     5,
		PageSize: 2,
	}
	
	// Mock the GetPaginated call for page beyond range
	expectedResponse := &repository.PaginationResponse[domain.Account]{
		Data:       []domain.Account{}, // Empty data for page beyond range
		Page:       5,
		PageSize:   2,
		Total:      3,
		TotalPages: 2,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 5 && req.PageSize == 2
	})).Return(expectedResponse, nil)
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	pagination := response.Pagination
	if pagination.Page != 5 {
		t.Errorf("Expected page 5, got %d", pagination.Page)
	}
	
	if pagination.Total != 3 {
		t.Errorf("Expected total 3, got %d", pagination.Total)
	}
	
	if pagination.TotalPages != 2 {
		t.Errorf("Expected total pages 2, got %d", pagination.TotalPages)
	}
	
	// Should return empty data for page beyond range
	if len(pagination.Data) != 0 {
		t.Errorf("Expected 0 accounts for page beyond range, got %d", len(pagination.Data))
	}
}

func TestGetAccountsHandler_Handle_ShouldRetrieveSingleAccountCorrectly(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(t)
	handler := NewGetAccountsHandler(mockRepo)
	
	// Create single test account
	account := domain.NewAccount("SINGLE-001", "Single User", domain.NewMoney(50000, domain.THB))
	account.ID = uuid.New()
	
	query := &queries.GetAccountsQuery{
		Page:     1,
		PageSize: 5,
	}
	
	// Mock the GetPaginated call for single account
	expectedResponse := &repository.PaginationResponse[domain.Account]{
		Data:       []domain.Account{*account},
		Page:       1,
		PageSize:   5,
		Total:      1,
		TotalPages: 1,
	}
	mockRepo.EXPECT().GetPaginated(mock.Anything, mock.MatchedBy(func(req repository.PaginationRequest) bool {
		return req.Page == 1 && req.PageSize == 5
	})).Return(expectedResponse, nil)
	
	ctx := context.Background()
	response, err := handler.Handle(ctx, query)
	
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
		t.Errorf("Expected 1 account, got %d", len(pagination.Data))
	}
	
	returnedAccount := pagination.Data[0]
	if returnedAccount.Number != "SINGLE-001" {
		t.Errorf("Expected account number 'SINGLE-001', got %s", returnedAccount.Number)
	}
	
	if returnedAccount.HolderName != "Single User" {
		t.Errorf("Expected holder name 'Single User', got %s", returnedAccount.HolderName)
	}
}