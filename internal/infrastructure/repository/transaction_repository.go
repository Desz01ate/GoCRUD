package repository

import (
	"arise_tech_assessment/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Repository[domain.Transaction, uuid.UUID]
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error)
	FindByAccountIDPaginated(ctx context.Context, accountID uuid.UUID, req PaginationRequest) (*PaginationResponse[domain.Transaction], error)
	FindByStatus(ctx context.Context, status domain.TransactionStatus) ([]domain.Transaction, error)
	FindByStatusPaginated(ctx context.Context, status domain.TransactionStatus, req PaginationRequest) (*PaginationResponse[domain.Transaction], error)
	FindByType(ctx context.Context, txType domain.TransactionType) ([]domain.Transaction, error)
	FindByTypePaginated(ctx context.Context, txType domain.TransactionType, req PaginationRequest) (*PaginationResponse[domain.Transaction], error)
	FindByReference(ctx context.Context, reference string) (*domain.Transaction, error)
	FindByDateRange(ctx context.Context, from, to time.Time) ([]domain.Transaction, error)
	FindByDateRangePaginated(ctx context.Context, from, to time.Time, req PaginationRequest) (*PaginationResponse[domain.Transaction], error)
}

type transactionRepository struct {
	*GormRepository[domain.Transaction, uuid.UUID]
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		GormRepository: NewGormRepository[domain.Transaction, uuid.UUID](db),
	}
}

func (r *transactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.GormRepository.db.WithContext(ctx).Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindByAccountIDPaginated(ctx context.Context, accountID uuid.UUID, req PaginationRequest) (*PaginationResponse[domain.Transaction], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var transactions []domain.Transaction
	var total int64

	query := r.GormRepository.db.WithContext(ctx).Where("from_account_id = ? OR to_account_id = ?", accountID, accountID)

	if err := query.Model(&domain.Transaction{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &PaginationResponse[domain.Transaction]{
		Data:       transactions,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *transactionRepository) FindByStatus(ctx context.Context, status domain.TransactionStatus) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.GormRepository.db.WithContext(ctx).Where("status = ?", status).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindByStatusPaginated(ctx context.Context, status domain.TransactionStatus, req PaginationRequest) (*PaginationResponse[domain.Transaction], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var transactions []domain.Transaction
	var total int64

	query := r.GormRepository.db.WithContext(ctx).Where("status = ?", status)

	if err := query.Model(&domain.Transaction{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &PaginationResponse[domain.Transaction]{
		Data:       transactions,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *transactionRepository) FindByType(ctx context.Context, txType domain.TransactionType) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.GormRepository.db.WithContext(ctx).Where("type = ?", txType).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindByTypePaginated(ctx context.Context, txType domain.TransactionType, req PaginationRequest) (*PaginationResponse[domain.Transaction], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var transactions []domain.Transaction
	var total int64

	query := r.GormRepository.db.WithContext(ctx).Where("type = ?", txType)

	if err := query.Model(&domain.Transaction{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &PaginationResponse[domain.Transaction]{
		Data:       transactions,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *transactionRepository) FindByReference(ctx context.Context, reference string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.GormRepository.db.WithContext(ctx).Where("reference = ?", reference).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByDateRange(ctx context.Context, from, to time.Time) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.GormRepository.db.WithContext(ctx).Where("created_at >= ? AND created_at <= ?", from, to).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindByDateRangePaginated(ctx context.Context, from, to time.Time, req PaginationRequest) (*PaginationResponse[domain.Transaction], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var transactions []domain.Transaction
	var total int64

	query := r.GormRepository.db.WithContext(ctx).Where("created_at >= ? AND created_at <= ?", from, to)

	if err := query.Model(&domain.Transaction{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &PaginationResponse[domain.Transaction]{
		Data:       transactions,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}
